Add-Type -AssemblyName System.Drawing
$sizes = 180, 192, 512
$outDir = Join-Path $PSScriptRoot "pwa"
foreach ($s in $sizes) {
  $bmp = New-Object System.Drawing.Bitmap $s, $s
  $g = [System.Drawing.Graphics]::FromImage($bmp)
  $g.Clear([System.Drawing.Color]::FromArgb(255, 26, 35, 50))
  $brush = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(255, 110, 168, 255))
  $fontSize = [Math]::Max(24, [int]($s / 5))
  $font = New-Object System.Drawing.Font("Segoe UI", $fontSize, [System.Drawing.FontStyle]::Bold)
  $g.DrawString("EW", $font, $brush, [int]($s * 0.2), [int]($s * 0.28))
  $g.Dispose()
  $path = Join-Path $outDir ("icon-$s.png")
  $bmp.Save($path, [System.Drawing.Imaging.ImageFormat]::Png)
  $bmp.Dispose()
}
