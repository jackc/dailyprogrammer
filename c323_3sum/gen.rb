# Generate test data for benchmarking implementations

500_000.times do
  puts (4 + rand(20)).times.reduce([]) { |accum, _| accum.push rand(-50..50) }.join(' ')
end
