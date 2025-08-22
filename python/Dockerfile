FROM --platform=$TARGETPLATFORM registry.cn-hangzhou.aliyuncs.com/kindlingx/autoinstrumentation-python:0.48b0 as python_agent

FROM python:3.11-slim

WORKDIR /app

ENV PYTHONPATH=/app
ENV PYTHONUNBUFFERED=1
ENV PORT=3500

RUN apt-get update && apt-get install -y \
    iproute2 \
    procps \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt


COPY . .

EXPOSE 3500

COPY --from=python_agent /autoinstrumentation /instrumentations/python
ENV PYTHONPATH /instrumentations/python:/instrumentations/python/opentelemetry/instrumentation/auto_instrumentation

CMD ["python", "main.py"]
