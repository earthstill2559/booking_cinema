# Cinema Booking App - Component Structure

## Overview
The Cinema Booking application has been refactored into reusable Vue 3 components for better maintainability, easier debugging, and improved UI/UX consistency.

## Component Architecture

### 📁 `/src/components/` Directory

#### 1. **Header.vue**
- **Purpose**: Application header with branding and title
- **Props**: None
- **Features**:
  - Professional gradient background (professional blue tones)
  - Logo with emoji icon
  - Title and subtitle
  - Responsive design for mobile/tablet/desktop
- **Styling**: Professional booking app aesthetic with blue primary color (#0066cc)

#### 2. **Screen.vue**
- **Purpose**: Cinema screen visualization
- **Props**: None
- **Features**:
  - SVG-based cinema screen display
  - Clean, minimalist design
  - Professional presentation
- **Styling**: Subtle shadow and responsive sizing

#### 3. **SeatsGrid.vue**
- **Purpose**: Seats selection grid display and interaction
- **Props**:
  - `seats` (Array, required): Array of seat objects with `id` and `booked` properties
  - `selectedSeats` (Array, required): Array of selected seat IDs
- **Emits**:
  - `toggle-seat(seatId)`: Emitted when a seat is clicked
- **Features**:
  - Professional grid layout with responsive columns
  - Three seat states: available (blue), selected (orange), booked (gray)
  - Hover effects and animations
  - Accessibility: disabled state for booked seats
  - Mobile-optimized touch interactions

#### 4. **Legend.vue**
- **Purpose**: Display seat status legend
- **Props**: None
- **Features**:
  - Visual legend showing available, selected, and booked seats
  - Responsive layout (horizontal on desktop, wraps on mobile)
  - Matches seat styling for visual consistency

#### 5. **BookingSummary.vue**
- **Purpose**: Display selected seats summary
- **Props**:
  - `selectedSeats` (Array, required): Array of selected seat IDs
- **Features**:
  - Shows selected seat IDs
  - Shows total count of selected seats
  - Smooth slide-down animation on appear
  - Only displays when seats are selected
  - Professional layout with orange accent border

#### 6. **Notification.vue**
- **Purpose**: Display success, error, and warning messages
- **Props**:
  - `message` (String, required): Message text to display
  - `messageType` (String, default: 'success'): Type of message ('success', 'error', 'warning')
- **Features**:
  - Color-coded notifications (green for success, red for error, orange for warning)
  - Smooth fade-slide animation
  - Professional styling with left accent border
  - Auto-dismiss support (parent component controls visibility)
  - Accessibility icons (✓, ✕, !)

#### 7. **BookingActions.vue**
- **Purpose**: Action buttons for booking confirmation and cancellation
- **Props**:
  - `selectedSeats` (Array, required): Array of selected seat IDs
  - `isLoading` (Boolean, default: false): Loading state indicator
- **Emits**:
  - `confirm-booking()`: Emitted when confirm button is clicked
  - `cancel-selection()`: Emitted when cancel button is clicked
- **Features**:
  - Primary (blue) confirm button
  - Secondary (white with border) cancel button
  - Loading spinner animation during booking
  - Disabled states with visual feedback
  - Responsive layout (side-by-side on desktop, stacked on mobile)

## Color Scheme

Professional ticket booking app styling inspired by major apps like Ticketmaster and BookMyShow:

- **Primary Blue**: `#0066cc` - Main action color
- **Primary Dark**: `#0052a3` - Hover/active states
- **Accent Orange**: `#ff6b35` - Highlight and selection
- **Accent Dark**: `#ee4266` - Alternative highlight
- **Success Green**: `#4caf50` - Success messages
- **Error Red**: `#f44336` - Error messages
- **Warning Orange**: `#ff9800` - Warning messages
- **Background**: `#f5f7fa` - Subtle light background
- **Text Dark**: `#333333` - Primary text
- **Text Light**: `#666666` - Secondary text
- **Border**: `#e0e0e0` - UI borders

## Key Improvements

### 1. **Modularity**
- Each component has a single responsibility
- Easy to update individual components without affecting others
- Components can be tested independently

### 2. **Maintainability**
- Clear component boundaries and props interface
- Easier to debug and fix issues in specific components
- Reusable components can be used in other projects

### 3. **UI/UX**
- Professional ticket booking app aesthetic
- Consistent color scheme across all components
- Smooth animations and transitions
- Better visual hierarchy
- Improved responsive design

### 4. **Developer Experience**
- Clear prop/emit interfaces
- Well-documented components
- Easy to extend with new features
- Better code organization

## Usage in App.vue

```vue
<template>
  <div class="app-container">
    <Header />
    
    <div class="main-content">
      <Screen />

      <SeatsGrid 
        :seats="seats" 
        :selectedSeats="selectedSeats"
        @toggle-seat="toggleSeatSelection"
      />

      <Legend />

      <BookingSummary :selectedSeats="selectedSeats" />

      <Notification 
        :message="message" 
        :messageType="messageType"
      />

      <BookingActions 
        :selectedSeats="selectedSeats"
        :isLoading="loading"
        @confirm-booking="confirmBooking"
        @cancel-selection="cancelSelection"
      />
    </div>
  </div>
</template>
```

## Responsive Breakpoints

All components follow these breakpoints:
- **Desktop**: >= 1024px
- **Tablet**: 768px - 1023px
- **Mobile**: < 768px
- **Small Mobile**: < 480px

Each component has tailored styles for each breakpoint for optimal user experience.

## Future Enhancements

Possible improvements for future iterations:
1. Add movie selection component
2. Add showtime selection component
3. Add price display component
4. Add payment methods component
5. Add booking history component
6. Add user profile component
7. Add favorites/wishlist component
8. Add seat sorting/filtering options
9. Add accessibility improvements (ARIA labels, keyboard navigation)
10. Add animations for seat selection feedback
