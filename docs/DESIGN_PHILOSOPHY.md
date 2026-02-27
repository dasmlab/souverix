# Design Philosophy & SDLC Methodology

## Foundation

Our design principles are based on fundamental SDLC practices that remain relevant across all modern ecosystems. We acknowledge the foundational work of:

- **Martin Fowler & Jez Humble** - "Continuous Integration" (the bible)
- **Per Andersson** - "You can just write that in Go easy..." - and we haven't looked back

## Core Principle

**Agile is not a design philosophy** - it's a way teams collectively use tools, gauges, lights, and indicators to agree on when things are good, bad, ugly, or "puppy" (production-ready). 

Production of software and hybrid platforms-as-code follow a **Development Method that should be respected regardless of how you self-organize**.

## The 6-Phase SDLC Framework

All phases can iterate on themselves, constantly generating updates for east and west phases (before and after). This is a continuous, iterative cycle.

### Phase 1: Use/Case Requirements

#### Technical Requirements
**"What is it that you want to build"**

- Tooling is not the point - it's just record keeping (numbered Excel sheet, JIRA, whatever - except Rational Rose, never that)
- The requirement is sitting and discussing against all different use cases, their nuances and variants
- What they should do, should not do, how they should do them, how they should be built
- While it's easy to get lost in technical "how-to", that is only one part
- All questions and answers come from the business question: **"What are we trying to solve, do, make better, improve, or maybe just prove or disprove - and how does it advance where we are today, to where we want to be tomorrow?"**

#### Business Requirements
**"Why is it that we want to build"**

- Technical strategic advantage
- Customer request
- Research grant
- For the hell of it cause it's "cool" (cool is costly 😉)
- The Business requirements must exist and drive the overall reason to be, exist - large or small within an org
- From the installation script to larger distributed suites, always relates somehow to the business
- Sometimes those relations are tough to see in large orgs with many moving parts, different temporal mechanisms and cycles
- The macro is always the same: **how does it push, contribute, gain, improve xyz aspect of our business?**
- **For any technical or product/program structural style requirements, if they are not traceable back to a business requirement, they are probably not priority**

#### Industrialization Requirements
These are program, product, organization style requirements, usually derived from Framework, ecosystem and similar in scaled and distributed organizations. With the advent of remote work and the shift of everything else as remote work becomes the norm, this is all companies.

- These requirements are where and how strategic, long-vision (north star, roadmap) terms that provide company-based business requirements come into play
- Vital as they ensure that the Business and Technical requirements are kept within an acceptable viable business framework
- **Yup, this is where budgets live** 😊

### Phase 2: Architecture

While many think this is just about Architects putting together stacks to meet requirements, **architecture also includes**:

- How you do your infra and surrounding setup
- Lifecycle management support
- Processes and regulations
- How interactions work
- Laying out the practical implementation aspects between units, tasks, and goals
- **Architecting the resources** in terms of:
  - People skills
  - Experience
  - Aptitudes and deficiencies
  - All that good stuff

All of this should be documented so that (with AI Agents today) it's really easy for everyone to know what process and things to build and create and how to follow. But more importantly, again, this should always explain **"why we are doing it this way or that"** from the point of view of how it solves business or industrialization requirements (security/technical etc. are in the native balliwick of the T in the TBI).

### Phase 3: Design

Here is where you put together the real SW/HW/cables and whatever (links, cards, boards, radio, shoes, whatever) and come up with some (there are things inside this we don't get into).

### Phase 4: Integration

Here is where you run whatever you are doing on whatever it's intended to do.

### Phase 5: Test the Hell Out of It

Keep doing it in a circle method. Continuous testing, continuous validation.

### Phase 6: Publish and Lifecycle

Patch and ship. Lifecycle management, updates, maintenance.

---

## Practical Build Flow

In our practical build, we want these happening:

### 1. Build
- Get a built container up in a registry that can be pulled down and started
- With our Unit test and our knowledge that there is always a diagnostic server/API available
- Activate the unit test pass/fail (since we are carrying our test code with us)

### 2. Unit Test
- Run the built with the target diagnostic and ensure it starts
- Passing unit diag/tests via diagnostic API calls

### 3. Publish (Trigger System Test)
- Publish to a repo and tag
- Trigger deployment to an environment suitable to run system test with other components around

### 4. Run System Test
- Take the new component under test against the backdrop of its overall other nodes and such
- Beat the shit out of it as appropriate for each component

### 5. Publish Result and Publish as Stable
- Run further tests (as we have outlined in our design MDs all in the projects) that do stability and others
- That will change the release (if passed) on published test to `-released`, `-validated` and other things
- We will get to that later

---

## Component CI/CD Template

This philosophy is implemented as a **boilerplate template** for all components:

```
Build → Unit Test (via Diagnostic API) → Publish → System Test → Publish Stable
```

Each component follows this flow, ensuring:
- Continuous integration
- Automated testing
- Diagnostic-driven validation
- System integration validation
- Stable release management

---

## Notes

- This methodology scales from single-person to 100-person teams
- It has been used and modified to fit most development setups
- The principles remain constant regardless of team size or organization structure
- Tooling is secondary - the methodology and philosophy are primary

---

*"You can just write that in Go easy..." - Per Andersson*
