import test_model as tm
from packaging import version
# pyright: basic

def test_parse_version():
    assert tm._parse_version("1.2.3") == (1,2,3)
    assert tm._parse_version("1.2.3") < (1,2,4)
    assert tm._parse_version("2.2.3") > (1,9,9)
    assert tm._parse_version("2.2.3") > (1,9)
    assert tm._parse_version("1.26.1rc1") > (1,26,1)
    version.parse
