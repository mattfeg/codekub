value = STDIN.gets.to_i
[200, 100, 50, 20, 10, 5, 2, 1].each do |note|
  puts "R$ #{note}: #{value / note}"
  value %= note
end
