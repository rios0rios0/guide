Scaffold a new frontend project following the 5-layer Clean Architecture (Domain, Service, Infrastructure, Presentation, Main).

For architecture details, refer to the Frontend Design section of the Architecture rule. For TypeScript standards, refer to the JavaScript rule. For testing patterns, refer to the Testing rule. For Makefile setup, refer to the CI/CD rule.

## 5-Layer Architecture

```
Domain            (most abstract -- entities and contracts)
  ^
Service           (API and external service communication)
  ^
Infrastructure    (technology-specific implementations)
  ^
Presentation      (views and UI components)
  ^
Main              (DI wiring, providers -- "dirty" layer)
```

**Dependencies ALWAYS point towards Domain.** Main is the only layer that knows about all others.

## Directory Structure

```
src/
├── domain/
│   ├── entities/
│   │   └── user.ts
│   └── contracts/
│       └── user_repository.ts
├── service/
│   └── user_service.ts
├── infrastructure/
│   └── http_user_repository.ts
├── presentation/
│   ├── components/
│   │   └── user_list.tsx
│   ├── pages/
│   │   └── users_page.tsx
│   └── hooks/
│       └── use_users.ts
├── main/
│   ├── providers/
│   │   └── user_provider.tsx
│   ├── di/
│   │   └── container.ts
│   └── app.tsx
└── index.tsx
```

## Step-by-Step

### 1. Create the Domain layer -- entities and repository contracts

```typescript
// src/domain/entities/user.ts
export interface User {
  id: string;
  name: string;
  email: string;
}

// src/domain/contracts/user_repository.ts
import { User } from '../entities/user';

export interface UserRepository {
  findAll(): Promise<User[]>;
  findById(id: string): Promise<User | null>;
}
```

### 2. Create the Service layer

```typescript
// src/service/user_service.ts
import { User } from '../domain/entities/user';
import { UserRepository } from '../domain/contracts/user_repository';

export class UserService {
  constructor(private readonly repository: UserRepository) {}

  async listUsers(): Promise<User[]> {
    return this.repository.findAll();
  }
}
```

### 3. Create the Infrastructure layer

```typescript
// src/infrastructure/http_user_repository.ts
import { User } from '../domain/entities/user';
import { UserRepository } from '../domain/contracts/user_repository';

export class HttpUserRepository implements UserRepository {
  async findAll(): Promise<User[]> {
    const response = await fetch('/api/users');
    return response.json();
  }

  async findById(id: string): Promise<User | null> {
    const response = await fetch(`/api/users/${id}`);
    return response.json();
  }
}
```

### 4. Create the Presentation layer -- hooks and components

```typescript
// src/presentation/hooks/use_users.ts
import { useState, useEffect } from 'react';
import { User } from '../../domain/entities/user';
import { UserService } from '../../service/user_service';

export const useUsers = (service: UserService) => {
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    service.listUsers().then(setUsers);
  }, [service]);

  return { users };
};
```

### 5. Create the Main layer -- DI wiring

```typescript
// src/main/di/container.ts
import { HttpUserRepository } from '../../infrastructure/http_user_repository';
import { UserService } from '../../service/user_service';

export const createContainer = () => {
  const userRepository = new HttpUserRepository();
  const userService = new UserService(userRepository);
  return { userService };
};
```

### 6. File naming

All filenames MUST use **snake_case**:

```
user_list.tsx       correct
UserList.tsx        wrong
```

### 7. TypeScript rules

- NEVER use `any` -- use `unknown` with type guards instead
- Add explicit type annotations on all function parameters and return types
- Prefer `const` over `let`, never use `var`
- Use optional chaining (`?.`) and nullish coalescing (`??`)

### 8. Create Makefile

The project Makefile must import from the shared [pipelines repository](https://github.com/rios0rios0/pipelines) and expose `lint`, `test`, and `sast` targets. Refer to the CI/CD rule for details.

## Testing

- Use **Jest** + **React Testing Library** (or Enzyme)
- Use **Faker** for generating test data
- Every test must use `// given`, `// when`, `// then` comment blocks
- Prefer `shallow` for unit tests, `mount` for integration tests
- Test: props rendering, user interactions, async operations, router push

Refer to the JavaScript rule (Testing section) and the Testing rule for full conventions.
