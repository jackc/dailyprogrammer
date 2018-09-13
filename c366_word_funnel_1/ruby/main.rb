def funnel(source, candidate)
  expected_skips = 1
  skips = 0

  source = source.each_char.to_a
  candidate = candidate.each_char.to_a

  if source.size - expected_skips != candidate.size
    return false
  end

  candidate_iter = candidate.each
  candidate_char = candidate_iter.next

  source.each do |s|
    if candidate_char == s
      candidate_char = begin
        candidate_iter.next
      rescue StopIteration
        nil
      end
    else
      skips += 1
      if skips > expected_skips
          return false
      end
    end
  end

  true
end

while gets
  words = $_.split
  if words.size != 2
    $stderr.puts "Invalid input. Requires exactly two words."
    next
  end

  puts "#{words[0]} #{words[1]} => #{funnel(words[0], words[1])}"
end
