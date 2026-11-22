require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

avg = ->(g) { (3 * g[0] + 3 * g[1] + 4 * g[2]) / 10 }
rng = Random.new(22)
cases = [[7.0, 8.0, 6.5], [10.0, 10.0, 10.0], [0.0, 0.0, 0.0], [0.0, 10.0, 5.0]]
while cases.size < 13
  g = Array.new(3) { rng.rand(0..100) / 10.0 }
  cases << g if safe?(avg.call(g), 1)
end
abort "tie in fixed case" unless cases.all? { |g| safe?(avg.call(g), 1) }
puts JSON.generate(cases.map { |g| g.map { |x| format("%.1f", x) }.join(" ") + "\n" })
