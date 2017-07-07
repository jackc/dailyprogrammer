# def f(word, word_idx, alphabet, remaining_slots)
#   alphabet.each do |letter|
#     word.push letter
#     if remaining_slots == 0
#       yield word
#     else
#       f(word, word_idx+1, alphabet, remaining_slots - 1)
#     end
#     word.pop
#   end

# end

# f("", 0, %w[+ - * /], 4) do |form|
#   puts form
# end



# 1 2 3
# 2 1 3
# 3 2 1
# 3 1 2
# 2 3 1
# 1 3 2

# + +
# + -
# + *
# + /
# - +
# - -
# - *
# - /
# * +
# * -
# * *
# * /
# / +
# / -
# / *
# / /

class RepPermIter
  include Enumerable
  attr_accessor :indices
  attr_accessor :cursor_pos
  attr_accessor :alphabet
  def initialize(alphabet:, size:)
    @alphabet = alphabet
    @indices = Array.new(size, 0)
    @indices[0] = -1
  end

  # [0 0 0 0]
  # [1 0 0 0]

  def each
    while nextperm
      yield indices.map { |i| alphabet[i] }
    end
  end

  def nextperm
    indices[0] += 1



    carry_idx = 0
    while carry_idx < indices.size && indices[carry_idx] == alphabet.size
      indices[carry_idx] = 0
      carry_idx += 1
      if carry_idx < indices.size
        indices[carry_idx] += 1
      end
    end

    carry_idx < indices.size
  end

end

puts RepPermIter.new(alphabet: %w[+ - * /], size: 2).count

