bill, percent, people = STDIN.read.split
bill = bill.to_f
service = bill * percent.to_i / 100
total = bill + service
puts format("Serviço: R$ %.2f", service)
puts format("Total: R$ %.2f", total)
puts format("Por pessoa: R$ %.2f", total / people.to_i)
