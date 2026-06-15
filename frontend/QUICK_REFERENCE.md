# Cinema Booking App - Quick Reference

## Component Hierarchy

```
App.vue (Main Component)
├── Header.vue
│   └── [Static - No props, No events]
│
├── Screen.vue
│   └── [Static - No props, No events]
│
├── SeatsGrid.vue
│   ├── Props:
│   │   ├── seats (Array of objects with id and booked)
│   │   └── selectedSeats (Array of seat IDs)
│   └── Events:
│       └── @toggle-seat(seatId)
│
├── Legend.vue
│   └── [Static - No props, No events]
│
├── BookingSummary.vue
│   ├── Props:
│   │   └── selectedSeats (Array)
│   └── [No events - Display only]
│
├── Notification.vue
│   ├── Props:
│   │   ├── message (String)
│   │   └── messageType (String: 'success' | 'error' | 'warning')
│   └── [No events - Display only]
│
└── BookingActions.vue
    ├── Props:
    │   ├── selectedSeats (Array)
    │   └── isLoading (Boolean)
    └── Events:
        ├── @confirm-booking()
        └── @cancel-selection()
```

## Data Flow

```
App.vue (Data Management)
├── seats (ref) ──────────→ SeatsGrid (display)
├── selectedSeats (ref) ───→ SeatsGrid + BookingSummary + BookingActions
├── loading (ref) ─────────→ BookingActions
├── message (ref) ──────────→ Notification
└── messageType (ref) ──────→ Notification

Events flowing back:
SeatsGrid ──toggle-seat──→ App (toggleSeatSelection method)
BookingActions ──confirm-booking──→ App (confirmBooking method)
BookingActions ──cancel-selection──→ App (cancelSelection method)
```

## Color Palette

| Element | Color Code | Usage |
|---------|-----------|-------|
| Primary Button | `#0066cc` | Confirm booking, primary actions |
| Primary Hover | `#0052a3` | Hover states on primary |
| Accent | `#ff6b35` | Selected seats, highlights |
| Accent Dark | `#ee4266` | Alternative highlight |
| Available Seat | `#0066cc` | Available for booking |
| Selected Seat | `#ff6b35` | User selected |
| Booked Seat | `#f5f5f5` | Already booked |
| Success | `#4caf50` | Success messages |
| Error | `#f44336` | Error messages |
| Warning | `#ff9800` | Warning messages |
| Background | `#f5f7fa` | App background |
| Text Primary | `#333333` | Main text |
| Text Secondary | `#666666` | Secondary text |
| Border | `#e0e0e0` | UI borders |

## State Management

```javascript
// In App.vue
const seats = ref([])                    // All available seats from API
const selectedSeats = ref([])            // Currently selected seat IDs
const loading = ref(false)               // Booking in progress
const message = ref("")                  // Notification message
const messageType = ref("success")       // Message type
```

## Event Flow for Booking

1. User clicks seat button
   - `SeatsGrid` emits `@toggle-seat` event
   - `App` receives and calls `toggleSeatSelection(seatId)`
   - `selectedSeats` array is updated

2. User clicks confirm button
   - `BookingActions` emits `@confirm-booking` event
   - `App` calls `confirmBooking()` method
   - API call to `/book/{seatId}`
   - On success: `message` and `messageType` are set
   - `selectedSeats` is cleared after confirmation

3. User clicks cancel button
   - `BookingActions` emits `@cancel-selection` event
   - `App` calls `cancelSelection()` method
   - `selectedSeats` array is cleared
   - `message` is cleared

## API Endpoints

```
GET /seats              → Returns array of all seats
POST /book/{seatId}     → Books a specific seat
```

## Styling Strategy

### Global Styles (`style.css`)
- CSS variables for colors
- Default typography
- Base element styles (body, html, links, etc.)
- Scrollbar styling
- Responsive breakpoints

### Component Styles (Individual `.vue` files)
- Component-specific layout
- Component-specific colors (using CSS variables)
- Component-specific animations
- Responsive adjustments for component

### Responsive Breakpoints

| Device | Breakpoint | Changes |
|--------|-----------|---------|
| Desktop | ≥ 1024px | Full layout, optimized spacing |
| Tablet | 768px - 1023px | Adjusted padding, smaller fonts |
| Mobile | < 768px | Stacked layout, small seat buttons |
| Small Mobile | < 480px | Minimal padding, extra small elements |

## Component Props Validation

All components use strict prop validation:

```javascript
defineProps({
  propName: {
    type: String,              // or Array, Boolean, Number, Object
    required: true,            // or false
    default: "default value",  // optional
    validator: (value) => {    // optional custom validator
      return /* validation logic */
    }
  }
})
```

## Common Patterns

### Toggle Selection
```javascript
const toggleSeatSelection = (seatId) => {
  const index = selectedSeats.value.indexOf(seatId)
  if (index > -1) {
    selectedSeats.value.splice(index, 1)
  } else {
    selectedSeats.value.push(seatId)
  }
}
```

### Async API Call
```javascript
const confirmBooking = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL}/book/${seatId}`, {
      method: "POST",
    })
    if (!res.ok) throw new Error("Failed")
    message.value = "Success!"
    messageType.value = "success"
  } catch (error) {
    message.value = "Error: " + error.message
    messageType.value = "error"
  } finally {
    loading.value = false
  }
}
```

### Conditional Rendering
```vue
<!-- Show only when condition is true -->
<div v-if="selectedSeats.length > 0">
  Selected: {{ selectedSeats.join(", ") }}
</div>

<!-- Transition animation -->
<transition name="slide-down">
  <div v-if="show">Content</div>
</transition>
```

### CSS Transitions
```css
.element {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.element:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}
```

## Debugging Checklist

- [ ] Component imports are correct
- [ ] Props are passed correctly
- [ ] Events are emitted with correct names
- [ ] Event listeners use correct names (@event-name)
- [ ] API URL is correct
- [ ] CSS classes match template
- [ ] No console errors
- [ ] Responsive design tested on all breakpoints
- [ ] All colors match design specification
- [ ] Animations are smooth
- [ ] Loading states work properly
- [ ] Error handling works properly

## Performance Tips

1. **Minimize Re-renders**: Use computed for derived state
2. **Debounce Events**: For frequent events like scroll/resize
3. **Lazy Load**: Large components loaded on demand
4. **Cache API Data**: Don't refetch same data repeatedly
5. **Optimize Assets**: Compress images, minimize CSS/JS
6. **Use Production Build**: For deployment
7. **Monitor Bundle Size**: Keep dependencies minimal

## Browser Support

- Chrome/Edge (latest 2 versions)
- Firefox (latest 2 versions)
- Safari (latest 2 versions)
- iOS Safari 12+
- Chrome Android (latest version)

## File Size Overview

```
Header.vue              ~1.5 KB
Screen.vue              ~1.2 KB
SeatsGrid.vue           ~2.8 KB
Legend.vue              ~1.8 KB
BookingSummary.vue      ~1.9 KB
Notification.vue        ~1.7 KB
BookingActions.vue      ~2.1 KB
App.vue                 ~3.2 KB
style.css               ~4.0 KB
────────────────────────────────
Total (uncompressed)   ~20 KB
```

## Next Steps

1. Test all components
2. Verify API integration
3. Test on different devices
4. Add error boundary component (optional)
5. Add loading skeleton (optional)
6. Add animation for seat selection (optional)
7. Add sorting/filtering features (optional)
