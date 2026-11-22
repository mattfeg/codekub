require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

ok = lambda do |c|
  s = c[0] * c[1] / 100
  safe?(s, 2) && safe?(c[0] + s, 2) && safe?((c[0] + s) / c[2], 2)
end
rng = Random.new(30)
cases = [[150.0, 10, 4], [89.9, 0, 3], [0.01, 20, 1], [10_000.0, 20, 50]]
cases.delete_if { |c| !ok.call(c) }
while cases.size < 13
  c = [rng.rand(1..1_000_000) / 100.0, rng.rand(0..20), rng.rand(1..50)]
  cases << c if ok.call(c)
end
puts JSON.generate(cases.map { |b, p, n| format("%.2f %d %d\n", b, p, n) })
