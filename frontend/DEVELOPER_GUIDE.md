# Cinema Booking App - Developer Guide

## Quick Start

### File Structure
```
frontend/
├── src/
│   ├── components/           # Reusable Vue components
│   │   ├── Header.vue       # App header with branding
│   │   ├── Screen.vue       # Cinema screen display
│   │   ├── SeatsGrid.vue    # Seat selection grid
│   │   ├── Legend.vue       # Seat status legend
│   │   ├── BookingSummary.vue    # Selected seats summary
│   │   ├── Notification.vue # Message notifications
│   │   └── BookingActions.vue    # Confirm/Cancel buttons
│   ├── App.vue              # Main app component
│   ├── main.js              # App entry point
│   └── style.css            # Global styles
├── package.json
└── COMPONENT_STRUCTURE.md   # Component documentation
```

## Common Tasks

### How to Fix Errors

1. **Notification Issues**: Edit `src/components/Notification.vue`
   - Message styling
   - Message type handling
   - Animation timing

2. **Seat Selection Issues**: Edit `src/components/SeatsGrid.vue`
   - Seat state handling
   - Click event handling
   - Grid layout responsiveness

3. **Booking Flow Issues**: Edit `src/App.vue`
   - API integration
   - State management
   - Event handling

4. **UI/Styling Issues**: 
   - Component-level styles: Edit individual `.vue` files
   - Global styles: Edit `src/style.css`

### How to Change Colors

All colors are defined in `src/style.css` using CSS variables:

```css
:root {
  --color-primary: #0066cc;      /* Main blue */
  --color-accent: #ff6b35;       /* Orange highlight */
  --color-success: #4caf50;      /* Green success */
  --color-error: #f44336;        /* Red error */
  --color-warning: #ff9800;      /* Orange warning */
  /* ... more colors ... */
}
```

To change colors:
1. Update the CSS variables in `src/style.css`
2. All components automatically use the new colors

### How to Change Layout

#### Seat Grid Columns
In `src/components/SeatsGrid.vue`:
```css
.seats-grid {
  grid-template-columns: repeat(auto-fit, minmax(60px, 1fr));
  gap: 14px;
}
```
Change `minmax(60px, 1fr)` to adjust column sizing.

#### Spacing and Padding
Each component controls its own spacing. Find the component and adjust padding/margin values.

#### Font Sizes
Modify individual component styles or global styles in `src/style.css`.

### How to Add New Features

#### Add a new message type
1. Update `Notification.vue`:
   ```vue
   :validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
   ```

2. Add styles in `<style>` section:
   ```css
   .notification.info {
     background: #e3f2fd;
     color: #1976d2;
     border-left-color: #2196f3;
   }
   ```

#### Add seat filtering
1. Create new component `src/components/SeatFilter.vue`
2. Import in `App.vue`
3. Add filter logic in `App.vue` script

#### Add seat sorting
1. Add sorting method in `App.vue`:
   ```javascript
   const sortSeats = (criteria) => {
     // sorting logic
   };
   ```
2. Create filter component and emit sorting events

## Testing Components

### Test Individual Components
Create a test file for each component:

```javascript
// BookingSummary.test.js
import { mount } from '@vue/test-utils'
import BookingSummary from '@/components/BookingSummary.vue'

describe('BookingSummary', () => {
  it('renders selected seats', () => {
    const wrapper = mount(BookingSummary, {
      props: {
        selectedSeats: ['A1', 'A2', 'A3']
      }
    })
    expect(wrapper.text()).toContain('A1, A2, A3')
  })
})
```

### Test App Flow
Test the complete booking flow in `App.test.js`:
- Fetch seats
- Select seats
- Confirm booking
- Check notifications

## Debugging Tips

1. **Vue DevTools**: Install Vue DevTools browser extension
   - Inspect component props and state
   - Track event emissions
   - Debug state changes in real-time

2. **Console Logging**: Add `console.log()` in methods
   ```javascript
   const toggleSeatSelection = (seatId) => {
     console.log('Toggling seat:', seatId)
     emit('toggle-seat', seatId)
   }
   ```

3. **Network Tab**: Check API requests
   - Verify correct endpoints
   - Check request/response data
   - Monitor loading states

4. **Component Inspection**: Right-click → Inspect Element
   - Check computed styles
   - Verify HTML structure
   - Check for CSS conflicts

## Performance Optimization

1. **Lazy Load Components** (if needed):
   ```javascript
   const Header = defineAsyncComponent(() => 
     import('./components/Header.vue')
   )
   ```

2. **Memoize Computed Properties**:
   ```javascript
   const availableSeats = computed(() => {
     return seats.value.filter(s => !s.booked)
   })
   ```

3. **Debounce Event Handlers**:
   ```javascript
   const debouncedToggleSeat = debounce(toggleSeatSelection, 300)
   ```

## Styling Best Practices

1. **Use CSS Variables**: Leverage predefined colors in `style.css`
2. **Mobile-First**: Always include mobile breakpoints
3. **Accessibility**: Ensure proper contrast ratios, hover states, disabled states
4. **Consistency**: Match padding, margins, and spacing across components
5. **Animations**: Use smooth transitions (0.2s - 0.4s recommended)

## Component Prop Types

All components use TypeScript-style prop definitions:
- `type`: Data type (Array, String, Boolean, Number)
- `required`: Is prop required?
- `default`: Default value if not provided
- `validator`: Custom validation function (optional)

Example:
```javascript
defineProps({
  message: {
    type: String,
    required: true,
  },
  messageType: {
    type: String,
    default: 'success',
    validator: (value) => ['success', 'error', 'warning'].includes(value),
  },
})
```

## Resources

- [Vue 3 Documentation](https://vuejs.org/)
- [Vue 3 Composition API](https://vuejs.org/guide/extras/composition-api-faq.html)
- [Vite Documentation](https://vitejs.dev/)
- [CSS Grid Guide](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Grid_Layout)
- [CSS Flexbox Guide](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout)

## Common Issues and Solutions

### Issue: Components not rendering
**Solution**: Check imports in App.vue, verify component file names match exactly

### Issue: Styles not applying
**Solution**: Check CSS specificity, clear browser cache, verify class names in template

### Issue: Events not firing
**Solution**: Check emit names match exactly, verify parent listening for correct event

### Issue: Props not updating
**Solution**: Ensure props are immutable, use computed for reactive values, check v-bind syntax

### Issue: API calls failing
**Solution**: Check backend is running, verify API_URL is correct, check CORS settings

## Contact & Support

For issues or questions:
1. Check this guide first
2. Review component comments
3. Check Vue documentation
4. Ask team members
