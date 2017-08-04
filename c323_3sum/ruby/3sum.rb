while gets
  nums = $_.split.map { |s| Integer(s) }.sort

  answers = []
  (0...nums.size).each do |i|
    ((i+1)...nums.size).each do |j|
      ((j+1)...nums.size).each do |k|
        answers.push [nums[i], nums[j], nums[k]] if nums[i] + nums[j] + nums[k] == 0
      end
    end
  end

  answers.sort.uniq.each do |a|
    puts a.join(" ")
  end

  puts
end
