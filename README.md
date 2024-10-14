# Highlights

It turns out, Entity-Componenent-System is basically just a relational database table with extreme constraints on lookup time.
As such, entities are indexed in an implemented-from-scratch [RoaringBitmap](https://arxiv.org/abs/1402.6407). We can make 100,000,000
inserts in 1 second. Nice!
