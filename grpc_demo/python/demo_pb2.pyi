from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class NumberStream(_message.Message):
    __slots__ = ["number"]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    number: int
    def __init__(self, number: _Optional[int] = ...) -> None: ...

class TestReply(_message.Message):
    __slots__ = ["counter", "message"]
    COUNTER_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    counter: int
    message: str
    def __init__(self, message: _Optional[str] = ..., counter: _Optional[int] = ...) -> None: ...

class TestRequest(_message.Message):
    __slots__ = ["counter", "message"]
    COUNTER_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    counter: int
    message: str
    def __init__(self, message: _Optional[str] = ..., counter: _Optional[int] = ...) -> None: ...
