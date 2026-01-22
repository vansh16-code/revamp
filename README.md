# Vehicle Rental Platform - Backend API

A comprehensive vehicle rental platform backend built with Go, Gin, and PostgreSQL. Enables students to rent out their vehicles when not in use and allows others to book them.

## Features

### 🔐 Authentication & Authorization
- JWT-based authentication
- Bcrypt password hashing
- Role-based access control (user/admin)
- Protected routes with middleware

### 👤 User Management
- Student profiles with course, department, year
- Dual roles: vehicle owner and renter
- Document verification (driving license, Aadhar, student ID)
- User statistics tracking

### 🚗 Vehicle Management
- Complete CRUD operations
- Vehicle types: bike, car, scooter, etc.
- Search and filter by type, price, location, rating
- OBD tracker integration
- Vehicle verification (RC, insurance, PUC)
- Image uploads and descriptions

### 📅 Availability Management
- Set time slots for vehicle availability
- Recurring availability support
- Conflict detection and overlap validation
- Real-time availability checking

### 📝 Booking System
- Create, confirm, cancel bookings
- Status flow: pending → confirmed → ongoing → completed
- Multiple pricing models:
  - Distance-based (per km)
  - Time-based (per hour)
  - Hybrid (distance + time)
- Security deposit handling
- Booking history and active bookings

### 🔢 OTP Verification
- **Pickup Flow**: Owner generates OTP → Renter verifies → Ride starts
- **Return Flow**: Owner generates OTP → Renter verifies → Ride completes
- Odometer and fuel level tracking
- Automatic calculations:
  - Distance traveled
  - Fuel consumed
  - Final price based on actual usage

### 📄 Document Verification
- Upload documents with URLs and metadata
- Document types:
  - User: Driving license, Student ID, Aadhar
  - Vehicle: RC, Insurance, PUC
- Admin approval/rejection workflow
- Automatic verification status updates
- Document expiry tracking

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL (Neon)
- **ORM**: GORM
- **Authentication**: JWT
- **Encryption**: AES-GCM for PII
- **Containerization**: Docker

## Project Structure

```
.
├── auth/               # JWT token generation and validation
├── config/             # Database configuration
├── handlers/           # HTTP request handlers
│   ├── user.go
│   ├── vehicle.go
│   ├── availability.go
│   ├── booking.go
│   └── document.go
├── middleware/         # HTTP middleware
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
├── models/             # Database models
│   ├── user.go
│   ├── vehicle.go
│   ├── availability.go
│   ├── booking.go
│   ├── document.go
│   ├── obd_tracker.go
│   └── constants.go
├── routes/             # Route definitions
├── utils/              # Utility functions
│   ├── encryption.go
│   ├── otp.go
│   └── pricing.go
├── main.go
├── Dockerfile
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker (optional)

### Environment Variables

Create a `.env` file in the root directory:

```env
DATABASE_URL=postgresql://user:password@host/database?sslmode=require
JWT_SECRET=your-secret-key-here
ENCRYPTION_KEY=your-32-byte-encryption-key
PORT=8080
```

### Installation

1. Clone the repository:
```bash
git clone https://github.com/vansh16-code/revamp.git
cd revamp
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

### Using Docker

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - Login user

### User Management
- `GET /api/profile` - Get current user profile
- `GET /api/users` - Get all users

### Vehicle Management
- `POST /api/vehicles` - Create vehicle
- `GET /api/vehicles` - Get all vehicles (with filters)
- `GET /api/vehicles/:id` - Get vehicle by ID
- `PUT /api/vehicles/:id` - Update vehicle
- `DELETE /api/vehicles/:id` - Delete vehicle
- `GET /api/my-vehicles` - Get my vehicles

### Availability Management
- `POST /api/vehicles/:id/availability` - Set availability
- `GET /api/vehicles/:id/availability` - Get vehicle availability
- `PUT /api/availability/:id` - Update availability
- `DELETE /api/availability/:id` - Delete availability
- `GET /api/availability/check` - Check availability

### Booking Management
- `POST /api/bookings` - Create booking
- `GET /api/bookings` - Get bookings (with filters)
- `GET /api/bookings/:id` - Get booking by ID
- `POST /api/bookings/:id/confirm` - Confirm booking (owner)
- `POST /api/bookings/:id/cancel` - Cancel booking
- `GET /api/bookings/active` - Get active booking
- `GET /api/bookings/history` - Get booking history

### OTP Verification
- `POST /api/bookings/:id/pickup/generate-otp` - Generate pickup OTP (owner)
- `POST /api/bookings/:id/pickup/verify-otp` - Verify pickup OTP (renter)
- `POST /api/bookings/:id/return/generate-otp` - Generate return OTP (owner)
- `POST /api/bookings/:id/return/verify-otp` - Verify return OTP (renter)

### Document Management
- `POST /api/documents` - Upload document
- `GET /api/documents` - Get my documents
- `GET /api/documents/:id` - Get document by ID
- `DELETE /api/documents/:id` - Delete document

### Admin Endpoints
- `GET /api/admin/documents/pending` - Get pending documents
- `POST /api/admin/documents/:id/verify` - Approve/reject document

## Database Models

### User
- Personal info (name, email, phone)
- Student details (ID, course, department, year)
- Documents (license, Aadhar)
- Ratings (as owner and renter)
- Verification status

### Vehicle
- Vehicle details (type, brand, model, year)
- Pricing (per km, per hour, per day, base price)
- Features (helmet, fuel type, transmission)
- Documents (RC, insurance, PUC)
- Location and availability
- Statistics (bookings, km driven)

### Booking
- Booking details (start/end time, location)
- Pricing model and calculations
- OTP verification data
- Odometer and fuel tracking
- Status tracking

### Document
- Document type and number
- Upload URL and metadata
- Verification status
- Admin approval/rejection

## Security Features

- JWT token authentication
- Password hashing with bcrypt
- AES-GCM encryption for PII (phone, student ID)
- Safe type assertions throughout
- Input validation and sanitization
- CORS middleware
- Protected admin routes

## Pricing Models

### Distance-Based
```
Final Price = Base Price + (Distance × Price Per Km)
```

### Time-Based
```
Final Price = Base Price + (Duration Hours × Price Per Hour)
```

### Hybrid
```
Final Price = Base Price + (Distance × Price Per Km) + (Duration Hours × Price Per Hour)
```

## Booking Lifecycle

1. **Create Booking** - Renter creates booking request
2. **Confirm Booking** - Owner confirms the booking
3. **Pickup** - Owner generates OTP → Renter verifies → Status: Ongoing
4. **Return** - Owner generates OTP → Renter verifies → Status: Completed
5. **Final Calculation** - System calculates final price based on actual usage

## Document Verification Flow

1. User uploads document with metadata
2. Document status: Pending
3. Admin reviews and approves/rejects
4. If approved, user/vehicle verification status updates
5. User fully verified when all required documents approved

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o app main.go
./app
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

Vansh - [@vansh16-code](https://github.com/vansh16-code)

Project Link: [https://github.com/vansh16-code/revamp](https://github.com/vansh16-code/revamp)
