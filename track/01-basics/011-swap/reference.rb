a, b = STDIN.read.split.map(&:to_i)
puts "antes: A = #{a}, B = #{b}"
tmp = a
a = b
b = tmp
puts "depois: A = #{a}, B = #{b}"
