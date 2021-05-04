from ctypes import CDLL, c_char_p
import os
from pathlib import Path

__all__ = ["ScanError", "scan"]
__version__ = "v0.17.2"


class ScanError(Exception):
    pass


def _errcheck(result: bytes, func, args):
    if result.startswith(b"ERROR"):
        msg = result.lstrip(b"ERROR ").decode()
        raise ScanError(msg)
    elif result == b"OK":
        return None
    else:
        try:
            return ScanError("Unexpected error: " + result.decode())
        except Exception:
            return ScanError("Unexpected error (can't decode error message)")


_default_loc = str(Path(__file__).parent.absolute() / "./libtrivy.so")
_path = os.getenv("LIBTRIVY_PATH", _default_loc)
_dll = CDLL(_path)
_dll.TrivyScan.argtypes = [c_char_p, c_char_p]
_dll.TrivyScan.restype = c_char_p
_dll.TrivyScan.errcheck = _errcheck


def scan(src: str, dst: str):
    _dll.TrivyScan(src.encode(), dst.encode())
