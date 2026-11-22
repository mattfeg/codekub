require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

ok = ->(c) { d = c[0] * c[1] / 100; safe?(d, 2) && safe?(c[0] - d, 2) }
rng = Random.new(26)
cases = [[200.0, 10], [59.9, 20], [100.0, 0], [100.0, 100], [0.01, 50]]
cases.delete_if { |c| !ok.call(c) }
while cases.size < 13
  c = [rng.rand(1..10_000_000) / 100.0, rng.rand(0..100)]
  cases << c if ok.call(c)
end
puts JSON.generate(cases.map { |p, d| format("%.2f %d\n", p, d) })
