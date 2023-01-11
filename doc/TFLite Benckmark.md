1. https://www.tensorflow.org/lite/performance/measurement
2. https://github.com/tensorflow/tensorflow/tree/master/tensorflow/lite/tools/benchmark#benchmark-multiple-performance-options-in-a-single-run


# 使用TFLite Benchmark tool

```bash
adb shell "mkdir -p /data/local/tmp/tflite_models"
adb push tmp/android_aarch64_benchmark_model_performance_options /data/local/tmp
adb shell "chmod +x /data/local/tmp/android_aarch64_benchmark_model_performance_options"
adb push model/mobilenet_quant_v1_224.tflite /data/local/tmp/tflite_models
adb shell "/data/local/tmp/android_aarch64_benchmark_model_performance_options \
  --num_threads=4 \
  --graph=/data/local/tmp/tflite_models/mobilenet_quant_v1_224.tflite \
  --warmup_runs=1 \
  --enable_op_profiling=true \
  --num_runs=10"
```