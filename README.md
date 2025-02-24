# Refactoring and Testing Legacy Code

This is a collection of examples and lessons around refactoring and testing legacy code. They're written in Go but designed to be language-agnostic.

I created this repo to supplement a talk I gave at an internal New York Times developer conference in 2023. Alas, I can't include the slide show I used, but rest assured it was amazing.

Each lesson has its own directory:

- [Lesson 1](./lesson1/) uses an example web app to demonstrate using dependency injection to unit test otherwise hard-to-test code.
- [Lesson 2](./lesson2/) covers how to find candidates for refactoring/testing in large codebases.
- [Lesson 3](./lesson3/) explains the dangers of relying too heavily on "code coverage" as a quality metric, and the importance of testing the right things.

## TLDR

- The best time to refactor/improve code is when you're already working in it. Don't wait for someone to ask you to do it. (But do check with your team if a refactor is going to take a long time!)
- Break code into logical pieces and use interfaces to inject dependencies. This makes unit testing much simpler
- If your refactor is significant, open a pull request _before_ you add any new functionality. Let reviewers focus on just the refactor.
- Refactoring code to make it more testable pays off quickly. Typically adding new functionality becomes much faster because you know exactly where to add it and how to test it.
- Practice test-driven development whenever it makes sense, but don't be dogmatic because sometimes it doesn't.