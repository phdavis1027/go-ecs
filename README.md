# Highlights (Papers and algorithms)

- It turns out, Entity-Componenent-System is basically just a relational database table with extreme constraints on lookup time.
As such, entities are indexed in an implemented-from-scratch [RoaringBitmap](https://arxiv.org/abs/1402.6407). We can make 100,000,000
inserts in 1 second. Nice!

- Update to the above: it appears that the 64-bit version is actually substantially more complicated than
the 32-bit version. In fact, a naive port is possible, but would take enormous amounts of storage. It seems real implementations leverage B-Trees
somehow, but I've already spent too much time on this data structure and I want to make progress on the actual game engine part, so I'm just
lifting a popular pre-built RoaringBitmap package.

- The ECS is looking like it'll become more central to the engine than I initially thought. 
That's because the system component provides a natural place to put a pretty (theoretically robust)
concurrency system. My plan is to implement [this paper](https://arxiv.org/pdf/1503.03642), where "transactions"/"transaction-pieces"
get mapped to Systems and "records" get mapped to entity types (which, in the ECS, are represented as 256 independent bitsets.) 
Since systems should have relatively independent concerns, most of the time, we should get a lot of concurrency from a scheme like this.

- Next step: I am *not* artistically gifted, so I have asked my sister to design a tileset. Which means that in the meantime, I'll be working on tile rendering. 
I think I'll use SDL for this application since I've played around with OpenGL before and spent most of my time fighting configuration rather than doing graphics
programming. Also, it seems that OpenGL integrates with SDL, so I can always come back and migrate piecemeal if I need hardware acceleration.

For this stage, the main event is probably going to be efficient rendering. Specifically, way back in elementary or moddle schol (before I knew anything about programming), I remember an episode
of Modern Marvels about how video players only render the pixels that change from frame to frame. For some reason I've always that idea was so cool, so I think it would be cool
to implement something like for myself.

- JK, I'm going to do it in OpenGL and learn incrementally. I really hate the feeling of not being able to look at code and understand what it's doing, and SDL
just felt to much like a blackbox. For example, I didn't like that it handled batching on its own. That feels anti-learning.
