"""Run once at build/install time; the runtime never downloads model weights."""
import argparse
import shutil
from pathlib import Path

from ultralytics import YOLO

parser = argparse.ArgumentParser()
parser.add_argument("--output", default="yolo26n.onnx")
parser.add_argument("--size", type=int, default=640)
args = parser.parse_args()
output = Path(args.output).resolve()
output.parent.mkdir(parents=True, exist_ok=True)
exported = Path(YOLO("yolo26n.pt").export(
    format="onnx", imgsz=args.size, batch=1, dynamic=False, simplify=False,
    opset=17, nms=False, device="cpu",
)).resolve()
if exported != output:
    shutil.copyfile(exported, output)
print(f"YOLO26n CPU model: {output}")
