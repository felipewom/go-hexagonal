# Golang Hexagonal Architeture study case

Project created to learn how to structure your golang/go programming language applications using the hexagonal architecture. #hexagonalarch 

## 📚Learning Objectives: 
    - Hexagonal Architecture/ Ports & Adapters
    - Go Programming
    - Restful
    - gRPC
    - Docker
    - Docker Compose
    - Dependency Injection (DI)
    - Inversion of Control (IoC)
    - Unit Tests
    - End To End Tests

___

## Hexagonal architecture (or Ports & Adapter pattern)
> proposed by [Alistair Cockburn](https://twitter.com/totheralistair) in 2005.


### Core
In this architecture, it is a technology agnostic component that contains all the business logic. The existence of other things in the outside should be completely ignored. In other words, the core shouldn’t be aware of how the application is served or where the the data is actually hold.

The core could be viewed as a “box” (represented as a hexagon) capable of resolve all the business logic independently of the infrastructure in which the application is mounted. This approach allow us to test the core in isolation and give us the ability to easily change infrastructure components.
Having a clear definition of the main component let us talk about the interactions with “things” that exists in the outside of the core. Those are what we call Actors.



### Actors
Actors are real world things that want to interact with the core. These things could be humans, databases or even other applications. Actors can be categorized into two groups, depending on who triggers the interaction:
Drivers (or primary) actors, are those who trigger the communication with the core. They do so to invoke a specific service on the core. A human or a CLI (command line interface) are perfect examples of drivers actors.
Driven (or secondary) actors, are those who are expecting the core to be the one who trigger the communication. In this case, is the core who needs something that the actor provides, so it sends a request to the actor and invoke a specific action on it. For example, if the core needs to save data into a MySQL database, then the core trigger the communication to execute an INSERT query on the MySQL client.

![diagram](https://miro.medium.com/max/1400/1*kEomMfgNPu1srEAH7-Z_LA.png)

Notice that the actors and the core “speak” different languages. An external application sends a request over http to perform a core service call (which does not understand what http means). Another example is when the core (which is technology agnostic) wants to save data into a mysql database (which speaks SQL).
Then, there must be “something” that can help us to make such translations. Here is where the Ports & Adapters come to play.


### Ports
In one hand, we have the ports which are interfaces that define how the communication between an actor and the core has to be done. Depending on the actor, the ports has different nature:
Ports for driver actors, define the set of actions that the core provides and expose to the outside. Each action generally correspond with a specific case of use.
Ports for driven actors, define the set of actions that the actor has to implement.
Notice that the ports belongs to the core. It is important, due to the core is the one who define which interactions are needed to achieve the business logic goals.

![diagram](https://miro.medium.com/max/550/1*b_c6bnop4qRjbK4ypUcWAg.png)

In black, the ports for driver actors. In gray, the ports for driven actors.


### Adapters
In the other hand, we have the adapters that are responsible of the transformation between a request from the actor to the core, and vice versa. This is necessary, because as we said earlier the actors and the core “speaks” different languages.
An adapter for a driver port, transforms a specific technology request into a call on a core service.
An adapter for a driven port, transforms a technology agnostic request from the core into an a specific technology request on the actor.

![diagram](https://miro.medium.com/max/1400/1*ERYx0IB1pN-5ZX98cKAoUw.png)

Actors connected to the core ports through adapters


### Dependency Injection
After the implementation is done, then it is necessary to connect, somehow, the adapters to the corresponding ports. This could be done when the application starts and it allow us to decide which adapter has to be connected in each port, this is what we call “Dependency injection”. For example, if we want to save data into a mysql database, then we just have to plug an adapter for a mysql database into the corresponding port or if we want to save data in memory (for testing) we need to plug an “in memory database” adapter into that port.

![diagram](https://miro.medium.com/max/1400/1*tXttGUY2PCCXW8CO6_Xg2w.png)
This image illustrate how easy is switching between infrastructure. Just change the dependency.