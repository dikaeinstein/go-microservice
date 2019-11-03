@search
Feature: As a user when I call the search endpoint, I would like to receive a list of kittens

  Scenario: User passes no search criteria
    Given I have no search criteria
    When I call the search endpoint
    Then I should receive a bad request message

  Scenario: User passes valid search criteria
    Given I have a valid search criteria
    When I call the search endpoint
    Then I should receive a list of kittens
