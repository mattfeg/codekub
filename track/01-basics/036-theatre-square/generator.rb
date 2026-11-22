require "json"

# Official Codeforces tests (CC BY 4.0, via open-r1/codeforces), in their original order.
puts JSON.generate(JSON.parse(File.read("official_tests.json")).map { |t| t["input"] })
