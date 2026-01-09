# TODO for CMS Dashboard and Register Updates

## UI Template Implementation

- [x] Update Login component to use MUI Sign-In template for attractive display.
- [x] Update Register component to use MUI Sign-Up template with grid layout and validation.
- [x] Update Dashboard component to use MUI Dashboard template with AppBar, collapsible Drawer, and grid layout.

## Error Handling Improvements

- [x] Add specific error messages for login (401: Invalid credentials).
- [x] Add specific error messages for register (409: User exists, generic for others).
- [x] Add detailed error handling in dashboard for API failures (401: Unauthorized with logout, others: generic).

## CTA and Routing

- [x] Ensure all buttons have complete CTAs with accurate routing (login to register, register back to login, logout to login).
- [x] Confirm workflow: register success -> login, login success -> dashboard.

## Notification Improvements

- [x] Add clear success notifications for login and register using MUI Snackbar.
- [x] Use detailed messages like "Login successful! Redirecting to dashboard..." and "Registration successful! Redirecting to login...".
- [x] Implement auto-redirect after success notifications.

## Testing and Verification

- [x] Run the application and verify all templates render correctly.
- [x] Verify error messages are specific and user-friendly.
- [x] Verify all CTAs and routing work as expected.
- [x] Verify Snackbar notifications appear for success and error cases.
