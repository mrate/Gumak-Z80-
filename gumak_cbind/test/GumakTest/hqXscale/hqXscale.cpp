#include <hqx/HQ2x.hh>
#include <hqx/HQ3x.hh>

extern "C" {
	__declspec(dllexport) void HQ3X(const void* input, uint32_t width, uint32_t height, void* output)
	{
		HQ3x h;
		h.resize(reinterpret_cast<const uint32_t*>(input), width, height, reinterpret_cast<uint32_t*>(output));
	}

	__declspec(dllexport) void HQ2X(const void* input, uint32_t width, uint32_t height, void* output)
	{
		HQ2x h;
		h.resize(reinterpret_cast<const uint32_t*>(input), width, height, reinterpret_cast<uint32_t*>(output));
	}
}