require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

ok = ->(c) { safe?(c * 9 / 5 + 32, 2) && safe?(c + 273.15, 2) }
rng = Random.new(21)
cs = [37.0, -40.0, 0.0, 100.0, -273.15, 10_000.0]
while cs.size < 13
  c = rng.rand(-27_315..1_000_000) / 100.0
  cs << c if ok.call(c)
end
abort "tie in fixed case" unless cs.all?(&ok)
puts JSON.generate(cs.map { |c| format("%.2f\n", c) })
