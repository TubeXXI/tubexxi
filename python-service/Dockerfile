# Use an official Python runtime as a parent image
FROM python:3.11-alpine

# Set the working directory in the container
WORKDIR /app

RUN pip install --no-cache-dir requests

# Copy the current directory contents into the container at /app
COPY . /app

RUN apt-get update && apt-get install -y --no-install-recommends \
		curl \
    wget \
    tar \
    gzip \
    unzip \
		chromium \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

COPY --from=builder /app/.env . 

# Make port 80 available to the world outside this container
EXPOSE 50051

# Run app.py when the container launches
CMD ["python", "app.py"]