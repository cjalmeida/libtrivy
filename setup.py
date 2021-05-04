import imp
from setuptools import setup, Extension
from setuptools.command.build_ext import build_ext as _build_ext
from wheel.bdist_wheel import bdist_wheel as _bdist_wheel
import subprocess
from distutils.util import get_platform
import os


class bdist_wheel(_bdist_wheel):
    # def finalize_options(self):
    #     super().finalize_options()
        # self.root_is_pure = False

    def get_tag(self):
        python, abi, plat = _bdist_wheel.get_tag(self)
        python, abi = "py3", "none"
        return python, abi, plat


setup(
    cmdclass={"bdist_wheel": bdist_wheel},
)
