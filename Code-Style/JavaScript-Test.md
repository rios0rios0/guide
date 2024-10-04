## Context
Front-end tests are usually quite challenging. Deciding __what__ and __how__ to test can sometimes lead to different opinions and perspectives.
Some other challenges are related to the way that any OpenSearch Dashboards plugin relates to some external packages, like the OpenSearch Dashboards, and ElasticSearch UI, both of which are very important.

This document should evolve with time. It's important to take the effort of maintaining it always up to date, especially when new modules and features require new tests techniques and tools.

## :hammer: Tools
We're using similar tools from the OSD tests suite.
1. Jest - [JavaScript Testing Framework](https://jestjs.io/)
2. Enzyme - [JavaScript Testing utility for React](https://enzymejs.github.io/enzyme/)
3. Faker - [Fake Data Generator](https://github.com/marak/Faker.js/)

## :black_joker: Unit Tests

### :dart: What To Test
At the unit tests, we're trying to test all of _our_ business rules and the behavior _our_ code. More than worrying if the view is being rendered (which is usually an elastic-ui responsibility), we are trying to ensure that all the possible _flows_ are leading to the desired _states_, and firing the proper _side effects_.

a. Test if the props are being passed properly (_our code_):
```ts
it('should render an ExampleChildren with the right props', () => {
  // given
  const component = mount(<ExampleComponent prop={'prop_value'} />);

  // then
  expect(component.find(ExampleChildren).prop('prop')).toBe('prop_value');
});
```

b. Test if a component has the right output based on the props:
```ts
it('should render the Content component', () => {
  // when
  const component = shallow(<ExampleComponent loading={false} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});

it('should render Loading component', () => {
  // when
  const component = shallow(<ExampleComponent loading={true} />);

  // then
  expect(component.find(Content).exists()).toBeFalsy();
  expect(component.find(Loading).exists()).toBeTruthy();
});
```

c. Test if a user interaction is properly mutating the app state:
```ts
it('should present the Congrats component when button clicked', async () => {
  // given
  const component = mount(<ExampleComponent />);

  // when
  component.find(ExampleButton).simulate('click');
  component.update();

  // then
  expect(component.find(Congrats).exists()).toBeTruthy();
});
```

d. Test if async actions are being executed properly:

```ts
import {flushComponent} from "tests/public/utils/flush_component.ts";

it('should render the Children with the loaded data', async () => {
  // given
  const exampleService = new MockExampleService();
  const exampleData = faker.random.alphaNumeric();
  exampleService.withSuccess(exampleData); // mocking the service response
  const component = mount( // mount() will call the useEffect hook, that in this case is loading data from the API.
    <ExampleAsyncComponent / >,
    {services: {exampleService}} // replace the context service with the mock
  );

  // For dealing with async updates React does you will find two different approachs in the code, please use the new one.
  // NEW
  await flushComponent(component);
  // OLD
  await act(async () => {
    await flushPromises(); // this function will force jest to wait all the promises to be resolved, in our case, the API requests.
  });
  component.update(); // this method will force our component to update. In this case, it already has the API data in the state.

  // then
  expect(component.find(Children).exists()).toBeTruthy();
  expect(component.find(Children).prop('loadedData')).toBeTruthy();
});
```

e. Test if a piece of logic returns the expected output:
Obs. Avoid repeating the business rules here, it's better to hard code the inputs and outputs to avoid bad tests.
```ts
it('should return the sum multiplied by two', () => {
  // given
  const numbersCollection = new NumbersCollection([1, 2, 3]); // don't make this 'dynamic', otherwise you'll have to duplicate the business rule here

  // when
  const result = numbersCollection.sumAndDuplicate(); // (1 + 2 + 3) * 2

  // then
  expect(result).toBe(12) // hardcode the expected response
});
```

f. Some components use the push method from useHistory hook in react-router-dom, test it like this:
```tsx
const mockHistoryPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'), // import the other methods as usual
  useHistory: () => ({
    push: mockHistoryPush, // mocks the push from useHistory
  }),
}));

it('should call push', async () => {
  // given
  const component = mount(<Something />);

  // when
  await flushComponent(component);

  // then
  expect(mockHistoryPush).toHaveBeenCalled();
});
```

### :question: How to test
We have some tools to help the tests happen. The `test_setup` file exports some of these functions.

a. `act` is a requirement for actions that perform state updates in the React components
b. `flushPromises` will force jest to wait for all the promises to be resolved, usually the API requests.
c. `mount` and `shallow` will render the components with all the wrappers applied. It's also possible to pass a partial context object, for mocking purposes.
i. The difference between the two is that `mount` will render the full DOM, while `shallow` will create a simpler representation of it. To see this in action, try to render the same component in both ways and run `console.log(component.debug())` in the test.
ii. `shallow` can be more performant since it doesn't render the full DOM.
iii. An important thing to notice is that `mount` will call the useEffect hook, while `shallow` will not.
