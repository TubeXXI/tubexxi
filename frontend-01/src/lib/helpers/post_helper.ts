import type { RequestEvent } from '@sveltejs/kit';
import { BaseHelper } from './base_helper';
import { PostSchema } from '@/utils/schema.js';
import readingTime from 'reading-time';
import matter from 'gray-matter';
import fs from 'fs/promises';
import path from 'path';
import { unified } from 'unified';
import remarkParse from 'remark-parse';
import { toc } from 'mdast-util-toc';

export class PostHelper extends BaseHelper {
	constructor(event: RequestEvent) {
		super(event);
	}

	/**
	 * Get all posts with pagination
	 * @param params - QueryParams
	 * @returns PaginatedResult<BlogPost>
	 */
	async getAllPosts(params: QueryParams = { page: 1, limit: 10 }): Promise<PaginatedResult<BlogPost>> {
		const {
			page = 1,
			limit = 10,
			search = '',
			sort_by = 'publishedDate',
			order_by = 'desc',
			status = 'published',
			tag,
			series,
			year,
			month
		} = params;

		const allPostFiles = import.meta.glob('/src/routes/content/**/*.md', { eager: true });
		const postEntries = Object.entries(allPostFiles);

		let posts: BlogPost[] = [];

		for (const [filePath, file] of postEntries) {
			const metadata = (file as any).metadata as PostSchema;

			if (!metadata || !metadata.slug) continue;

			const pathParts = filePath.split('/');
			const slugFromPath = pathParts[pathParts.length - 1].replace('.md', '');

			const finalSlug = metadata.slug || slugFromPath;

			const raw = await fs.readFile(
				path.resolve('src/routes/content', `${finalSlug}.md`),
				'utf-8'
			);

			const { content, data } = matter(raw);
			const stats = readingTime(content);


			posts.push({
				meta: {
					...metadata,
					slug: finalSlug
				},
				path: filePath.replace('/src/routes/content/', '').replace('.md', ''),
				readingTime: Math.ceil(stats.minutes),
				words: stats.words,
				headings: this.extractHeadings(content)
			});
		}


		let filteredPosts = posts;

		if (status && status !== 'ALL' && status !== 'all') {
			filteredPosts = filteredPosts.filter(post =>
				post.meta.status?.toLowerCase() === status.toLowerCase()
			);
		}

		if (search) {
			const searchLower = search.toLowerCase();
			filteredPosts = filteredPosts.filter(post =>
				post.meta.title?.toLowerCase().includes(searchLower) ||
				post.meta.description?.toLowerCase().includes(searchLower) ||
				post.meta.tags?.some(tag => tag.toLowerCase().includes(searchLower)) ||
				false
			);
		}

		if (tag) {
			filteredPosts = filteredPosts.filter(post =>
				post.meta.tags?.includes(tag)
			);
		}

		if (series) {
			filteredPosts = filteredPosts.filter(post =>
				post.meta.series?.title === this.normalizeParams({ series })[series]
			);
		}

		if (year) {
			filteredPosts = filteredPosts.filter(post => {
				if (!post.meta.publishedDate) return false;
				const postYear = new Date(post.meta.publishedDate).getFullYear();
				return postYear === year;
			});
		}

		if (year && month) {
			filteredPosts = filteredPosts.filter(post => {
				if (!post.meta.publishedDate) return false;
				const postDate = new Date(post.meta.publishedDate);
				return postDate.getFullYear() === year &&
					(postDate.getMonth() + 1) === month;
			});
		}

		filteredPosts.sort((a, b) => {
			const aValue = (a.meta as any)[sort_by];
			const bValue = (b.meta as any)[sort_by];

			if (sort_by === 'publishedDate' || sort_by === 'lastUpdatedDate') {
				const aDate = new Date(aValue as string || 0).getTime();
				const bDate = new Date(bValue as string || 0).getTime();
				return order_by === 'asc' ? aDate - bDate : bDate - aDate;
			}

			if (typeof aValue === 'string' && typeof bValue === 'string') {
				return order_by === 'asc'
					? aValue.localeCompare(bValue)
					: bValue.localeCompare(aValue);
			}

			if (typeof aValue === 'number' && typeof bValue === 'number') {
				return order_by === 'asc' ? aValue - bValue : bValue - aValue;
			}

			return 0;
		});

		const totalItems = filteredPosts.length;
		const totalPages = Math.ceil(totalItems / limit);
		const currentPage = Math.min(Math.max(page, 1), totalPages || 1);
		const startIndex = (currentPage - 1) * limit;
		const endIndex = startIndex + limit;

		const paginatedData = filteredPosts.slice(startIndex, endIndex);

		return {
			data: paginatedData,
			pagination: {
				current_page: currentPage,
				total_pages: totalPages,
				total_items: totalItems,
				has_next: currentPage < totalPages,
				has_prev: currentPage > 1,
				limit: limit,
				next_page: currentPage < totalPages ? currentPage + 1 : undefined,
				prev_page: currentPage > 1 ? currentPage - 1 : undefined
			}
		};
	}

	/**
	 * Get a post by its slug
	 * @param slug - The slug of the post
	 * @returns BlogPost | null
	 */
	async getPostBySlug(slug: string): Promise<BlogPost | null> {
		const raw = await fs.readFile(
			path.resolve('src/routes/content', `${slug}.md`),
			'utf-8'
		);

		const { data } = matter(raw);

		const parsed = PostSchema.parse(data);

		const allPosts = await this.getAllPosts({ page: 1, limit: 100 });
		const currentPost = allPosts.data.find(post => post.meta.slug === slug);

		return {
			meta: parsed,
			// @ts-expect-error
			component: currentPost?.default,
			readingTime: currentPost?.readingTime || 0,
			words: currentPost?.words || 0,
			path: currentPost?.path || '',
			headings: currentPost?.headings || []
		};
	}
	/**
	 * Get related posts based on tags and series
	 * @param currentSlug - The slug of the current post
	 * @param limit - The number of related posts to return (default: 3)
	 * @returns BlogPost[]
	 */
	async getRelatedPosts(currentSlug: string, limit: number = 3): Promise<BlogPost[]> {
		const allPosts = await this.getAllPosts({ page: 1, limit: 100 });
		const currentPost = allPosts.data.find(post => post.meta.slug === currentSlug);

		if (!currentPost) return [];

		const posts = allPosts.data
			.filter(post => post.meta.slug !== currentSlug)
			.filter(post => {
				// Match by tags
				const commonTags = currentPost.meta.tags?.filter(tag =>
					post.meta.tags?.includes(tag)
				) || [];

				// Match by series
				const sameSeries = currentPost.meta.series?.title === post.meta.series?.title;

				return commonTags.length > 0 || sameSeries;
			})
			.slice(0, limit);

		return posts.map(post => ({
			...post,
			readingTime: post.readingTime,
			words: post.words,
			headings: post.headings
		}));
	}

	/**
	 * Get all unique tags
	 * @returns string[]
	 */
	async getAllTags(): Promise<string[]> {
		const allPosts = await this.getAllPosts({ page: 1, limit: 1000 });
		const tags = new Set<string>();

		allPosts.data.forEach(post => {
			post.meta.tags?.forEach(tag => tags.add(tag));
		});

		return Array.from(tags).sort();
	}

	/**
	 * Get all unique series
	 * @returns string[]
	 */
	async getAllSeries(): Promise<string[]> {
		const allPosts = await this.getAllPosts({ page: 1, limit: 1000 });
		const series = new Set<string>();

		allPosts.data.forEach(post => {
			if (post.meta.series?.title) {
				series.add(post.meta.series.title);
			}
		});

		return Array.from(series).sort();
	}

	/**
	 * Get posts by year
	 * @param year - The year to filter posts by
	 * @returns BlogPost[]
	 */
	async getPostsByYear(year: number): Promise<BlogPost[]> {
		const allPosts = await this.getAllPosts({ page: 1, limit: 1000 });
		const posts = allPosts.data.filter(post => {
			const postYear = new Date(post.meta.publishedDate).getFullYear();
			return postYear === year;
		});

		return posts.map(post => ({
			...post,
			readingTime: post.readingTime,
			words: post.words,
			headings: post.headings
		}));
	}

	/**
	 * Get archive data (years with post counts)
	 * @returns Array<{ year: number; count: number }>
	 */
	async getArchiveData(): Promise<Array<{ year: number; count: number }>> {
		const allPosts = await this.getAllPosts({ page: 1, limit: 1000 });
		const archiveMap = new Map<number, number>();

		allPosts.data.forEach(post => {
			const year = new Date(post.meta.publishedDate).getFullYear();
			archiveMap.set(year, (archiveMap.get(year) || 0) + 1);
		});

		return Array.from(archiveMap.entries())
			.map(([year, count]) => ({ year, count }))
			.sort((a, b) => b.year - a.year);
	}

	/**
	 * Get all unique slug names
	 * @returns string[]
	 */
	async getAllSlugNames(): Promise<string[]> {
		const allPosts = await this.getAllPosts({ page: 1, limit: 1000 });
		const slugs = new Set<string>();

		allPosts.data.forEach(post => {
			slugs.add(post.meta.slug);
		});

		return Array.from(slugs).sort();
	}


	/**
	 * Normalize query parameters to handle single and array values uniformly
	 * @param params - The query parameters to normalize
	 * @returns Record<string, string | string[]>
	 */
	private normalizeParams(params: Record<string, string | string[] | undefined>) {
		const normalized: Record<string, string | string[]> = {};

		for (const [key, value] of Object.entries(params)) {
			if (value === undefined || value === '') continue;
			normalized[key] = Array.isArray(value) ? value : [value];
		}

		return normalized;
	}

	/**
	 * Extract headings from markdown content
	 * @param markdown - The markdown content to extract headings from
	 * @returns Array<{ depth: number; value: string; slug: string }>
	 */
	private extractHeadings(markdown: string) {
		const tree = unified().use(remarkParse).parse(markdown);

		const result = toc(tree, {
			heading: null,
			maxDepth: 3
		});

		if (!result.map) return [];

		const headings: {
			value: string;
			depth: number;
			slug: string;
		}[] = [];

		function walk(node: any, depth = 0) {
			if (node.type === 'listItem') {
				const textNode = node.children?.[0]?.children?.[0];
				const text = textNode?.value;

				const link = node.children?.[0]?.children?.[1];
				const slug = link?.url?.replace('#', '');

				if (text && slug) {
					headings.push({
						value: text,
						depth,
						slug
					});
				}
			}

			if (node.children) {
				node.children.forEach((child: any) => walk(child, depth + 1));
			}
		}

		walk(result.map);

		return headings;
	}
}
