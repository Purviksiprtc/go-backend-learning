## Paota vs Raw RabbitMQ

### Code Complexity
- Raw RabbitMQ requires manual channel, exchange, queue, binding, ack handling
- Paota abstracts all AMQP boilerplate via WorkerPool

### Setup & Configuration
- Raw RabbitMQ needs explicit declarations
- Paota uses env-based config and sane defaults

### Reliability
- Raw requires manual retry & DLQ wiring
- Paota provides retry, delay, DLQ out of the box

### When to use Paota
- Rapid development
- Standard async workflows
- Clean architecture

### When to use Raw RabbitMQ
- Extremely fine-grained AMQP control
- Custom broker behaviors
