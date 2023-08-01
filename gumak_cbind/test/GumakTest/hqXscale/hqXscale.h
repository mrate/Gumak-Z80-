#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

	extern __declspec(dllexport) void HQ3X(const void* input, uint32_t width, uint32_t height, void* output);
	extern __declspec(dllexport) void HQ2X(const void* input, uint32_t width, uint32_t height, void* output);

#ifdef __cplusplus
}
#endif