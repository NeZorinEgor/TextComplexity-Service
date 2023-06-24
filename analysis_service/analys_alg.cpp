#include "analys_alg.h"

//getAverageSentenceLength | getAverageWordsLength | getConcret;
enum class Params{
	Sentence = 0,
	Word = 1,
	Concret = 2
};
static constexpr int kSizeOfParameters = 3;

namespace aal {

Result analys(const std::string& text) {
	using namespace Impl;

	Result result_out;

	auto [sen_length,word_length] = getAverageSentenceWordLength(text);
	auto concrent = getConcret(text);

//	WATER
	ParCombinator<int> water_combinator(kSizeOfParameters);
	auto & sentense_ranges = water_combinator[static_cast<int>(Params::Sentence)];
		sentense_ranges.push_range(1, 0.0, 5.0);
		sentense_ranges.push_range(2, 5.0, 6.0);
		sentense_ranges.push_range(3, 6.0, 7.0);
		sentense_ranges.push_range(4, 7.0, 8.0);
		sentense_ranges.push_range(5, 8.0, 9.0);
		sentense_ranges.push_range(6, 10.0, 11.0);
		sentense_ranges.push_range(7, 12.0, 13.0);
		sentense_ranges.push_range(8, 14.0, 16.0);
		sentense_ranges.push_range(9, 18.0, 22.0);
		sentense_ranges.push_range(10, 22.0, 10000.0);
	water_combinator.add_value(static_cast<int>(Params::Sentence), sen_length);

	auto&  word_ranges = water_combinator[static_cast<int>(Params::Word)];
		word_ranges.push_range(1, 0.0, 5.0);
		word_ranges.push_range(2, 5.0, 5.5);
		word_ranges.push_range(3, 5.5, 6.0);
		word_ranges.push_range(4, 6.0, 6.5);
		word_ranges.push_range(5, 6.5, 7.0);
		word_ranges.push_range(6, 7.0, 7.8);
		word_ranges.push_range(7, 7.8, 8.5);
		word_ranges.push_range(8, 8.5, 9.3);
		word_ranges.push_range(9, 9.3, 10.0);
		word_ranges.push_range(10, 10.0, 10000.0);
	water_combinator.add_value(static_cast<int>(Params::Word), word_length);

	auto& concret_ranges = water_combinator[static_cast<int>(Params::Concret)];
		concret_ranges.push_range(1, 0.0, 0.001);
		concret_ranges.push_range(2, 0.001, 0.01);
		concret_ranges.push_range(3, 0.01, 0.1);
		concret_ranges.push_range(4, 0.1, 1);
		concret_ranges.push_range(5, 1, 3);
		concret_ranges.push_range(6, 3, 6);
		concret_ranges.push_range(7, 6, 9);
		concret_ranges.push_range(8, 9, 12);
		concret_ranges.push_range(9, 12, 16);
		concret_ranges.push_range(10, 16, 10000);
	water_combinator.add_value(static_cast<int>(Params::Concret), concrent);

	result_out.water = water_combinator.getResult();
//	MOOD
	using tm = Result::TypeMood;

	ParCombinator<tm> mood_combinator(kSizeOfParameters);
	auto& sentense_ranges_mood = mood_combinator[static_cast<int>(Params::Sentence)];
		sentense_ranges_mood.push_range(tm::Sad, 6.0, 7.5);
		sentense_ranges_mood.push_range(tm::Happy, 4.5, 5.5);
		sentense_ranges_mood.push_range(tm::Lovely, 7.5, 8.5);
		sentense_ranges_mood.push_range(tm::Terrible, 9.0, 11.0);
		sentense_ranges_mood.push_range(tm::Boring, 10.0, 12.0);
	mood_combinator.add_value(static_cast<int>(Params::Sentence), sen_length);

	auto& word_ranges_mood = mood_combinator[static_cast<int>(Params::Word)];
		word_ranges_mood.push_range(tm::Sad, 9.0, 10.0);
		word_ranges_mood.push_range(tm::Happy, 8.0, 9.0);
		word_ranges_mood.push_range(tm::Lovely, 6.0, 7.0);
		word_ranges_mood.push_range(tm::Terrible, 9.0, 10.0);
		word_ranges_mood.push_range(tm::Boring, 7.0, 11.0);
	mood_combinator.add_value(static_cast<int>(Params::Word), word_length);

	auto& concret_ranges_mood = mood_combinator[static_cast<int>(Params::Concret)];
		concret_ranges_mood.push_range(tm::Sad, 0, 0);
		concret_ranges_mood.push_range(tm::Happy, 0, 0.5);
		concret_ranges_mood.push_range(tm::Lovely, 0, 0);
		concret_ranges_mood.push_range(tm::Terrible, 0, 0.5);
		concret_ranges_mood.push_range(tm::Boring, 0.5, 1000);
	mood_combinator.add_value(static_cast<int>(Params::Concret), concrent);

	result_out.mood = mood_combinator.getResult();
//	HARD
	ParCombinator<int> hard_combinator(kSizeOfParameters);
	auto& sentense_ranges_hard = hard_combinator[static_cast<int>(Params::Sentence)];
		sentense_ranges_hard.push_range(1, 0.0, 5.0);
		sentense_ranges_hard.push_range(2, 5.0, 6.0);
		sentense_ranges_hard.push_range(3, 6.0, 7.0);
		sentense_ranges_hard.push_range(4, 7.0, 8.0);
		sentense_ranges_hard.push_range(5, 8.0, 9.0);
		sentense_ranges_hard.push_range(6, 10.0, 11.0);
		sentense_ranges_hard.push_range(7, 12.0, 13.0);
		sentense_ranges_hard.push_range(8, 14.0, 16.0);
		sentense_ranges_hard.push_range(9, 18.0, 22.0);
		sentense_ranges_hard.push_range(10, 22.0, 10000.0);
	hard_combinator.add_value(static_cast<int>(Params::Sentence), sen_length);

	auto& word_ranges_hard = hard_combinator[static_cast<int>(Params::Word)];
		word_ranges_hard.push_range(1, 11.0, 12.0);
		word_ranges_hard.push_range(2, 10.5, 11.0);
		word_ranges_hard.push_range(3, 10.0, 10.5);
		word_ranges_hard.push_range(4, 9.0, 10.0);
		word_ranges_hard.push_range(5, 8.0, 9.0);
		word_ranges_hard.push_range(6, 7.0, 8.0);
		word_ranges_hard.push_range(7, 6.5, 7.5);
		word_ranges_hard.push_range(8, 6.0, 6.5);
		word_ranges_hard.push_range(9, 5.5, 6.0);
		word_ranges_hard.push_range(10, 0.0, 5.5);
	hard_combinator.add_value(static_cast<int>(Params::Word), word_length);

	auto& concret_ranges_hard = hard_combinator[static_cast<int>(Params::Concret)];
		concret_ranges_hard.push_range(1, 0.0, 0.005);
		concret_ranges_hard.push_range(2, 0.005, 0.05);
		concret_ranges_hard.push_range(3, 0.05, 0.5);
		concret_ranges_hard.push_range(4, 0.5, 1.5);
		concret_ranges_hard.push_range(5, 1.5, 3.5);
		concret_ranges_hard.push_range(6, 3.5, 6.5);
		concret_ranges_hard.push_range(7, 6.5, 9.5);
		concret_ranges_hard.push_range(8, 9.5, 12.5);
		concret_ranges_hard.push_range(9, 12.5, 16.5);
		concret_ranges_hard.push_range(10, 16.5, 10000);
	hard_combinator.add_value(static_cast<int>(Params::Concret), concrent);

	result_out.hard = hard_combinator.getResult();

	return result_out;
}

namespace Impl {

std::pair<float,float> getAverageSentenceWordLength(const std::string& text) {
	auto last_pos = text.begin();
	auto new_pos = text.begin();
	uint64_t word_count = 0, sum = 0;
	while (last_pos != text.end()) {
		new_pos = std::find_if( last_pos + 1, text.end(), [&](char ch) {
			if (ch == ' ' || ch == '\n' || ch == '\t')
				word_count++;
			return ch == '.' || ch == '?' || ch == '!'; //optimization
		});
		sum++;
		last_pos = new_pos;
	}
	return { float(word_count) / sum, (text.size() - word_count) / float(word_count) };
}


float getConcret(const std::string& text) {
	auto count_spaces = std::count_if(text.begin(), text.end(), [](char ch) { return ch == ' ' || ch == '\n' || ch == '\t'; });
	return float(std::count_if(text.begin(), text.end(), [](char ch) { return ch >= '0' && ch <= '9'; })) / ((text.size() - count_spaces) / 1000);
}

}

}