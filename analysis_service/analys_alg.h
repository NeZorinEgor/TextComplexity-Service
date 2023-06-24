#pragma once

#include <string>
#include <map>
#include <vector>
#include <memory>
#include <algorithm>

namespace aal {

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

Result analys(const std::string&);

namespace Impl {

std::pair<float, float> getAverageSentenceWordLength(const std::string&);
//отношение чисел и слов в ковычках ко всем другим 
//количество чисел на 1000 слов
float getConcret(const std::string&);

template <class T>
class ParCombinator {
public:
	class ParHolder {

		friend class ParCombinator;

		struct rangeContainer {
			float begin, end;
			T value;
		};

	public:

		void push_range(T value, float begin, float end) {
			m_ranges.push_back({ begin,end,value });
			if(m_counter->find(value) == m_counter->end())
				m_counter->operator[](value) = 0;
		}

	private:
		void set_counter(std::map<T, int>* mp) {
			m_counter = mp;
		}

		void procces_adding(float value) {
			for (auto& range : m_ranges) {
				if (range.begin <= value && range.end >= value) {
					m_counter->operator[](range.value)++;
				}
			}
		}
		std::map<T, int>* m_counter;
		std::vector<rangeContainer> m_ranges;
	};

	ParCombinator(int size) {
		m_holders.reset(new std::vector <ParHolder>(size));
		for (auto& par : *m_holders) {
			par.set_counter(&m_counter);
		}
	}

	ParHolder& operator[](int index) {
		return m_holders->at(index);
	}

	bool add_value(int row, float value) {
		if (m_holders->size() <= row)
			return false;
		m_holders->operator[](row).procces_adding(value);
		return true;
	}

	T getResult() {
		if constexpr (std::is_same_v<T, int>) {
			int sum_out = 0;
			int count_out = 0;
			for (auto [key, value] : m_counter) {
				sum_out += key * value;
				count_out += value;
			}
			return (count_out ? sum_out / count_out : 0); //TODO COUNT
		}

		if constexpr (std::is_same_v<T, Result::TypeMood>) {
			auto best_mood = std::max_element(m_counter.begin(), m_counter.end(),
				[](const std::pair<T, int>& p1, const std::pair<T, int>& p2) {
					return p1.second < p2.second; });
			return best_mood->first;
		}
	}

private:
	std::unique_ptr<std::vector<ParHolder>> m_holders;
	std::map<T, int> m_counter;
};

}

}

