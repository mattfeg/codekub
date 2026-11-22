a, b = STDIN.read.split.map(&:to_i)
q = a / b
r = a % b
puts "quociente: #{q}"
puts "resto: #{r}"
puts "#{a} = #{q} x #{b} + #{r}"
