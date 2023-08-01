// GumakTest.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <vector>
#include <chrono>
#include <thread>
#include <string_view>

#define __SIZE_TYPE__ size_t
#define _Complex

#include "../../../lib/gumak.h"

#define HQ3X

#ifdef HQ3X
#include <hqx/HQ2x.hh>
#include <hqx/HQ3x.hh>

#define STB_IMAGE_WRITE_IMPLEMENTATION
#include <stb_image_write.h>
#endif

struct Timer {
	Timer(const char* name) : name_{ name }, start_{ clock::now() } {}
	~Timer() {
		const auto elapsed{ std::chrono::duration_cast<std::chrono::microseconds>(clock::now() - start_).count() };
		std::cout << "  " << name_ << " took " << static_cast<float>(elapsed) / 1000.0f << "ms\n";
	}

	using clock = std::chrono::high_resolution_clock;
	const char* name_;
	clock::time_point start_;
};

struct AccuTimer {
	AccuTimer(const char* name, uint64_t& count, float& time) : name_{ name }, start_{ clock::now() }, count_{ count }, time_{ time } {}
	~AccuTimer() {
		const auto elapsed{ static_cast<float>(std::chrono::duration_cast<std::chrono::microseconds>(clock::now() - start_).count()) / 1000.0f };
		++count_;
		time_ = (count_ * time_ + elapsed) / (count_ + 1);
	}

	using clock = std::chrono::high_resolution_clock;
	const char* name_;
	clock::time_point start_;

	uint64_t& count_;
	float& time_;
};

void GenerateGrid(int width, int height) {
	std::vector<uint32_t> pixels(width * height);

	for (int x{}; x < width; ++x) {
		for (int y{}; y < height; ++y) {
			const auto grid{ (x % 8) == 0 || (y % 8) == 0 };

			pixels[y * width + x] = grid ? 0xff000000 : 0x00000000;
		}
	}

	stbi_write_png("grid.png", width, height, 4, pixels.data(), width * 4);
}

int main()
{
	std::cout << "Hello World!\n";

	GoString romPath{};
	romPath.p = "..\\..\\gumak\\roms";
	romPath.n = strlen(romPath.p);

	const auto inst{ GumakCreate(0, 48000, 1024, romPath) };
	if (inst == 0) {
		std::cerr << "Failed to create instance\n";
		return 0;
	}

	GoString snapshot{};
	snapshot.p = "..\\..\\gumak\\roms\\tmp.z80";
	snapshot.n = strlen(snapshot.p);

	const auto res{ GumakResolution() };
	const auto width{ static_cast<uint32_t>(res.r0) };
	const auto height{ static_cast<uint32_t>(res.r1) };

	GenerateGrid(width + 1, height + 1);

	std::vector<uint32_t> pixels(width * height);

#ifdef HQ3X
	const auto upscaledWidth{ static_cast<uint32_t>(3 * res.r0) };
	const auto upscaledHeight{ static_cast<uint32_t>(3 * res.r1) };

	std::vector<uint32_t> upscaled(upscaledWidth * upscaledHeight);

	HQ2x hq2;
	HQ3x hq3;

	uint64_t hq3xCount{};
	float hq3xTime{};
#endif

	GoSlice data;
	data.data = pixels.data();
	data.cap = data.len = pixels.size() * 4;

	int counter{};
	auto screenshot{ true };

	uint64_t displayCount{};
	float displayTime{};

	std::vector<uint8_t> buffers[2];
	buffers[0].resize(1024);
	buffers[1].resize(1024);
	GumakAudioBuffer(inst, 0, buffers[0].data(), 1024);
	GumakAudioBuffer(inst, 1, buffers[1].data(), 1024);

	while (true) {
		GumakUpdateFrame(inst);
		{
			AccuTimer t{ "GukamDisplayData", displayCount, displayTime };
			GumakDisplayDataRGB(inst, data, width * 4);
		}

#ifdef HQ3X
		{
			AccuTimer t{ "h3qx", hq3xCount, hq3xTime };
			hq3.resize(pixels.data(), width, height, upscaled.data());
		}
#endif

		switch (++counter) {
		case 80:
			std::cout << "Loading snapshot '" << snapshot.p << "'\n";

			if (GumakLoadSnapshot(inst, snapshot) == 0) {
				std::cout << "Failed to load snapshot.\n";
			}
			GumakHandleKey(inst, GumakKeyEnter, 1);
			break;
		case 90:
			GumakHandleKey(inst, GumakKeyEnter, 0);
			break;
		case 100:
			std::cout << "Arrow down pressed\n";

			GumakHandleKey(inst, GumakKeyShift, 1);
			GumakHandleKey(inst, GumakKey6, 1);
			break;
		case 101:
			std::cout << "Arrow down released\n";

			GumakHandleKey(inst, GumakKeyShift, 0);
			GumakHandleKey(inst, GumakKey6, 0);
			break;
		case 150:
			if (screenshot) {
				screenshot = false;

				std::cout << "=> Capture screenshot 'original.png'\n";
				stbi_write_png("original.png", width, height, 4, pixels.data(), width * 4);
#ifdef HQ3X
				{
					std::cout << "=> Capture screenshot 'resized_3x.png'\n";
					stbi_write_png("resized_3x.png", upscaledWidth, upscaledHeight, 4, upscaled.data(), upscaledWidth * 4);

					std::cout << "HQ2X upscale\n";
					{
						Timer t{ "h2qx" };
						hq2.resize(pixels.data(), width, height, upscaled.data());
					}

					std::cout << "=> Capture screenshot 'resized_2x.png'\n";
					stbi_write_png("resized_2x.png", width * 2, height * 2, 4, upscaled.data(), width * 2 * 4);
				}
#endif
			}

			counter = 0;
			break;
		}

		if (counter == 0) {
			break;
		}
	}

	std::cout << "Destroy instance\n";
	GumakDestroy(inst);

	std::cout << "GumakDisplayData: " << displayCount << "x - avg. " << displayTime << "ms\n";

#ifdef HQ3X
	std::cout << "HQ3X: " << hq3xCount << "x - avg. " << hq3xTime << "ms\n";

	std::cout << "=> Upscaled frame: avg. " << (hq3xTime + displayTime) << "ms\n";
#endif

	std::cout << "Done.\n";
	return 0;
}
