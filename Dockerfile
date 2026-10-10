FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
COPY static ./static
COPY detection ./detection
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/neonvr .

# PyTorch/Ultralytics are used only to export the official weights at build time.
FROM python:3.12-slim-bookworm AS model
ENV YOLO_AUTOINSTALL=False OMP_NUM_THREADS=2
RUN apt-get update && apt-get install -y --no-install-recommends libgl1 libglib2.0-0 && rm -rf /var/lib/apt/lists/*
WORKDIR /model
RUN pip install --no-cache-dir torch==2.9.1 torchvision==0.24.1 --index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir ultralytics==8.4.175 onnx==1.23.2
COPY detection/export_model.py ./
RUN python export_model.py --output /model/yolo26n.onnx

# CPU ONNX Runtime supports the Intel i7-6700; no CUDA/GPU runtime is installed.
FROM python:3.12-slim-bookworm
RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
COPY detection/requirements.txt /tmp/detection-requirements.txt
RUN pip install --no-cache-dir -r /tmp/detection-requirements.txt && rm /tmp/detection-requirements.txt
COPY --from=model /model/yolo26n.onnx /opt/neonvr/models/yolo26n.onnx
ENV NEONVR_MODEL_DIR=/opt/neonvr/models OMP_NUM_THREADS=2
COPY --from=build /out/neonvr /usr/local/bin/neonvr
EXPOSE 5500
ENTRYPOINT ["/usr/local/bin/neonvr"]
CMD ["-config", "/config/neonvr.json"]
