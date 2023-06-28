#pragma once

struct Result {
	enum class TypeMood {
		Sad = 0,
		Happy = 1,
		Lovely = 2,
		Terrible = 3,
		Boring = 4
	};

	int water;
	TypeMood mood;
	int hard;
};