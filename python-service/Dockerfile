FROM python:3.11-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1 \
	PYTHONUNBUFFERED=1

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
		chromium \
		chromium-driver \
		fonts-liberation \
		libasound2 \
		libatk-bridge2.0-0 \
		libatk1.0-0 \
		libc6 \
		libdrm2 \
		libgbm1 \
		libgtk-3-0 \
		libnss3 \
		libx11-xcb1 \
		libxcomposite1 \
		libxdamage1 \
		libxrandr2 \
		libu2f-udev \
		xdg-utils \
	&& rm -rf /var/lib/apt/lists/*

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . ./

EXPOSE 50051

CMD ["python", "server.py"]
