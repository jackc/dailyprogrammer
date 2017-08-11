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

RepPermIter.new(alphabet: %w[1 2 3 4 4 5], size: 3).each do |perm|
  p perm
end


