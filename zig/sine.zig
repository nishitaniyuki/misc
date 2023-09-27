const c = @cImport(@cInclude("soundio/soundio.h"));
const std = @import("std");

fn sio_err(err: c_int) !void {
    switch (err) {
        c.SoundIoErrorNone => {},
        c.SoundIoErrorNoMem => return error.NoMem,
        c.SoundIoErrorInitAudioBackend => return error.InitAudioBackend,
        c.SoundIoErrorSystemResources => return error.SystemResources,
        c.SoundIoErrorOpeningDevice => return error.OpeningDevice,
        c.SoundIoErrorNoSuchDevice => return error.NoSuchDevice,
        c.SoundIoErrorInvalid => return error.Invalid,
        c.SoundIoErrorBackendUnavailable => return error.BackendUnavailable,
        c.SoundIoErrorStreaming => return error.Streaming,
        c.SoundIoErrorIncompatibleDevice => return error.IncompatibleDevice,
        c.SoundIoErrorNoSuchClient => return error.NoSuchClient,
        c.SoundIoErrorIncompatibleBackend => return error.IncompatibleBackend,
        c.SoundIoErrorBackendDisconnected => return error.BackendDisconnected,
        c.SoundIoErrorInterrupted => return error.Interrupted,
        c.SoundIoErrorUnderflow => return error.Underflow,
        c.SoundIoErrorEncodingString => return error.EncodingString,
        else => return error.Unknown,
    }
}
