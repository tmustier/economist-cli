class EconomistTui < Formula
  desc "Terminal UI to browse and read The Economist"
  homepage "https://github.com/tmustier/economist-tui"
  head "https://github.com/tmustier/economist-tui.git"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"economist", "-ldflags", "-s -w", "."
  end

  test do
    system "#{bin}/economist", "--version"
  end

  def caveats
    <<~EOS
      Requires Chrome or Chromium for login and article fetching.
    EOS
  end
end
