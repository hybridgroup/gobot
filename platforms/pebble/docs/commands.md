# Functions

## PublishEvent(name string, data string)

It publishes an event.

#### Params

- **name** - **string** - event name
- **data** - **string** - value

#### API Command

**PublishEventC**

## PendingMessage()

It returns messages to be sent as notifications to pebble (Not recommended to be used directly)

#### API Command

**PendingMessageC**

## SendNotification(message string)

Sends notification to watch.

#### Params

- **message** - **string** - notification text

#### API Command

**SendNotification**
