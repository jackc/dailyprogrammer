while gets
  $_.split.map { |s| Integer(s) }
    .sort
    .combination(3).select do |a, b, c|
      a + b + c == 0
    end.sort.uniq.each do |a|
      puts a.join(" ")
    end

  puts
end
