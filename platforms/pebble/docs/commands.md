# Functions

## PublishEvent(name string, data string)

It publishes an event.

#### Params

- **name** - **string** - event name
- **data** - **string** - value

#### API Command

**publish_event**

## PendingMessage()

It returns messages to be sent as notifications to pebble (Not intented to be used directly)

#### API Command

**pending_message**

## SendNotification(message string)

Sends notification to watch.

#### Params

- **message** - **string** - notification text

#### API Command

**send_notification**
