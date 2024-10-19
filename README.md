# Highlights (Papers and algorithms)

- It turns out, Entity-Componenent-System is basically just a relational database table with extreme constraints on lookup time.
As such, entities are indexed in an implemented-from-scratch [RoaringBitmap](https://arxiv.org/abs/1402.6407). We can make 100,000,000
inserts in 1 second. Nice!

- The ECS is looking like it'll become more central to the system than I initially thought. 
That's because the system component provides a natural place to put a pretty (theoretically robusts)
concurrency system. My plan is to implement [this paper](https://arxiv.org/pdf/1503.03642), which "transactions"/"transaction-pieces"
get mapped to Systems and "records" get mapped to entity types (which, in the ECS, are represented as 256 independent bitsets. 
Since systems should have relatively independent concerns, most of the time, should get a lot of concurrency from a scheme like this.
