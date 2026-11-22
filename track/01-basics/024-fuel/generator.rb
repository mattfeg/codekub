require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

ok = ->(c) { safe?(c[0] / c[1], 3) && safe?(c[1] * c[2] / c[0], 2) }
rng = Random.new(24)
cases = [[500, 40.0, 6.0], [123, 9.87, 5.49], [1, 0.01, 0.01], [100_000, 10_000.0, 20.0]]
while cases.size < 13
  c = [rng.rand(1..100_000), rng.rand(1..1_000_000) / 100.0, rng.rand(1..2000) / 100.0]
  cases << c if ok.call(c)
end
abort "tie in fixed case" unless cases.all?(&ok)
puts JSON.generate(cases.map { |k, l, p| format("%d\n%.2f\n%.2f\n", k, l, p) })
