# Plugin Life Cycle

> **TL;DR:** OSD plugins have two lifecycle methods: `setup` (define routes, UI elements, and APIs) and `start` (initialize the plugin, register event handlers, and interact with other plugins).

## Overview

OpenSearch Dashboards plugins follow a defined lifecycle with two key methods that are called during initialization.

## `setup` Method

The `setup` method is called during plugin initialization. It defines the plugin's behavior and its interaction with the OSD platform:

- Register HTTP routes and API endpoints.
- Define UI modules and components.
- Configure notifications and platform integrations.

```ts
import { CoreSetup } from 'opensearch-dashboards/src/core/public';

export const setup = (core: CoreSetup): void => {
  const { http, uiModules, notifications } = core;

  http.get('/api/my_route', async (req, res) => {
    // Handle API request
  });

  uiModules.registerModule('myModule', {
    // Define UI components
  });

  notifications.toasts.addSuccess({
    title: 'My Plugin',
    text: 'Setup complete!',
  });
};
```

## `start` Method

The `start` method is called after `setup` completes for all plugins. It performs runtime initialization:

- Register event handlers and start background tasks.
- Interact with other OSD plugins and services.
- Access data sources and search indices.

```ts
import { CoreStart } from 'opensearch-dashboards/src/core/public';

export const start = (core: CoreStart): void => {
  const { savedObjects, data } = core;

  savedObjects.get('myDataSource')
    .then(dataSource => {
      // Use the data source
    });

  data.search({
    index: 'myIndex',
    body: {
      query: {
        match_all: {},
      },
    },
  })
    .then(res => {
      // Handle search results
    });
};
```

## Summary

| Method  | Timing                       | Purpose                                                               |
|---------|------------------------------|-----------------------------------------------------------------------|
| `setup` | During initialization        | Define routes, UI elements, APIs, and platform integrations           |
| `start` | After all plugins are set up | Initialize runtime behavior, interact with other plugins and services |

## References

- [Kibana Plugin API](https://www.elastic.co/guide/en/kibana/master/kibana-platform-plugin-api.html)
- [OpenSearch Dashboards Developer Guide](https://github.com/opensearch-project/OpenSearch-Dashboards/blob/main/DEVELOPER_GUIDE.md)
