#import <Cocoa/Cocoa.h>
#import <math.h>

extern void goAnimalsDesktopTick(void);
extern void goAnimalsDesktopKeyDown(void);
extern void goAnimalsDesktopSetSceneWidth(int width);
extern void goAnimalsDesktopSetSpeed(int speed);
extern void goAnimalsDesktopSetPetCount(int count);
extern void goAnimalsDesktopSetWheelEnabled(int enabled);
extern void goAnimalsDesktopSetMode(int mode);
extern void goAnimalsDesktopSetLanguage(int language);
extern void goAnimalsDesktopSetDisplayID(long long displayID);
extern void goAnimalsDesktopSetCoatMode(int mode);
extern void goAnimalsDesktopSetVariant(int variant);
extern void goAnimalsDesktopSetSelectedCoat(int index, int variant);
extern void goAnimalsDesktopSetNameLabels(int enabled);
extern void goAnimalsDesktopSetPetName(int index, char *value);
extern void goAnimalsDesktopSetPetSize(int index, int percent);
extern int goAnimalsDesktopClick(int x, int y);
extern int goAnimalsDesktopPetAt(int x, int y);
extern int goAnimalsDesktopGetSpeed(void);
extern int goAnimalsDesktopGetPetCount(void);
extern int goAnimalsDesktopGetWheelEnabled(void);
extern int goAnimalsDesktopGetMode(void);
extern int goAnimalsDesktopGetLanguage(void);
extern long long goAnimalsDesktopGetDisplayID(void);
extern int goAnimalsDesktopGetCoatMode(void);
extern int goAnimalsDesktopGetVariant(void);
extern int goAnimalsDesktopGetSelectedCoat(int index);
extern int goAnimalsDesktopGetVariantCount(void);
extern int goAnimalsDesktopCopyVariantLabel(int index, char *buffer, int length);
extern int goAnimalsDesktopCopyVariantGroupLabel(int index, char *buffer, int length);
extern int goAnimalsDesktopGetPetSize(int index);
extern int goAnimalsDesktopGetNameLabels(void);
extern int goAnimalsDesktopCopyPetName(int index, char *buffer, int length);
extern int goAnimalsDesktopGetPetDrawX(int index);
extern int goAnimalsDesktopGetPetDrawY(int index);

enum {
	AnimalsMenuSettings = 1001,
	AnimalsMenuSpeedSlow = 1101,
	AnimalsMenuSpeedNormal = 1103,
	AnimalsMenuSpeedFast = 1105,
	AnimalsMenuCountBase = 1200,
	AnimalsMenuWheelEnabled = 1301,
	AnimalsMenuLanguageJA = 1401,
	AnimalsMenuLanguageEN = 1402,
	AnimalsMenuCoatFixed = 1501,
	AnimalsMenuCoatSelected = 1502,
	AnimalsMenuCoatRandom = 1503,
	AnimalsMenuModeKeyboard = 1601,
	AnimalsMenuModeRandom = 1602,
	AnimalsMenuDisplayBase = 1700,
	AnimalsMenuVariantBase = 2000,
};

static const NSInteger AnimalsMaxPetCount = 10;
static const CGFloat AnimalsSpriteWidth = 96.0;
static const NSInteger AnimalsMinPetSize = 70;
static const NSInteger AnimalsMaxPetSize = 120;
static const NSInteger AnimalsPetSizeStep = 10;

static NSString *AnimalsVariantLabel(NSInteger index) {
	char buffer[256] = {0};
	int copied = goAnimalsDesktopCopyVariantLabel((int)index, buffer, (int)sizeof(buffer));
	if (copied <= 0) {
		return [NSString stringWithFormat:@"Animal %ld", (long)index + 1];
	}
	NSString *label = [NSString stringWithUTF8String:buffer];
	return label != nil ? label : [NSString stringWithFormat:@"Animal %ld", (long)index + 1];
}

static NSString *AnimalsVariantGroupLabel(NSInteger index) {
	char buffer[256] = {0};
	int copied = goAnimalsDesktopCopyVariantGroupLabel((int)index, buffer, (int)sizeof(buffer));
	if (copied <= 0) {
		return @"Animals";
	}
	NSString *label = [NSString stringWithUTF8String:buffer];
	return label != nil ? label : @"Animals";
}

static NSString *AnimalsVariantDisplayLabel(NSInteger index) {
	return [NSString stringWithFormat:@"%@ / %@", AnimalsVariantGroupLabel(index), AnimalsVariantLabel(index)];
}

static NSString *AnimalsPetSizeLabel(NSInteger percent) {
	return [NSString stringWithFormat:@"%ld%%", (long)percent];
}

static long long AnimalsScreenID(NSScreen *screen) {
	NSNumber *number = [[screen deviceDescription] objectForKey:@"NSScreenNumber"];
	return number != nil ? [number longLongValue] : 0;
}

static NSScreen *AnimalsScreenForDisplayID(long long displayID) {
	NSArray *screens = [NSScreen screens];
	if ([screens count] == 0) {
		return [NSScreen mainScreen];
	}
	if (displayID > 0) {
		for (NSScreen *screen in screens) {
			if (AnimalsScreenID(screen) == displayID) {
				return screen;
			}
		}
	}
	NSScreen *main = [NSScreen mainScreen];
	return main != nil ? main : [screens objectAtIndex:0];
}

static long long AnimalsActiveDisplayID(void) {
	long long displayID = goAnimalsDesktopGetDisplayID();
	return AnimalsScreenID(AnimalsScreenForDisplayID(displayID));
}

static NSString *AnimalsText(NSString *key) {
	BOOL english = goAnimalsDesktopGetLanguage() == 1;
	NSDictionary *ja = @{
		@"settingsOpen": @"設定を開く...",
		@"speed": @"速さ",
		@"speedSlow": @"ゆっくり",
		@"speedNormal": @"ふつう",
		@"speedFast": @"はやい",
		@"petCount": @"表示数",
		@"petCountUnit": @"匹",
		@"keyboardReaction": @"キーボード反応",
		@"language": @"Language",
		@"display": @"表示先ディスプレイ",
		@"quit": @"終了",
		@"settingsTitle": @"Animals Desktop 設定",
		@"support": @"対応OS: macOS 12 Monterey 以降 / Intel・Apple Silicon",
		@"tabAnimals": @"動物",
		@"tabMotion": @"動き",
		@"tabNames": @"名前",
		@"visibleCount": @"出現数",
		@"animalMode": @"動物の決め方",
		@"coatFixed": @"固定",
		@"coatSelected": @"1匹ずつ選ぶ",
		@"coatRandom": @"ランダム",
		@"fixedAnimal": @"固定する動物",
		@"perPetAnimal": @"1匹ごとの動物",
		@"petIndexSuffix": @"匹目",
		@"mode": @"モード",
		@"modeKeyboard": @"キーボード反応",
		@"modeRandom": @"ランダム散歩",
		@"speedLabel": @"速度",
		@"typingWheel": @"チンチラ/ハムスター回し車",
		@"petSize": @"サイズ",
		@"nameLabels": @"名前を表示",
		@"nameHint": @"ONのとき、動物にカーソルを乗せると名前が表示されます。",
		@"defaultPetName": @"どうぶつ",
		@"macNote": @"Mac版はメニューバー常駐です。Dockアイコンは通常表示しません。",
		@"close": @"閉じる",
	};
	NSDictionary *en = @{
		@"settingsOpen": @"Open Settings...",
		@"speed": @"Speed",
		@"speedSlow": @"Slow",
		@"speedNormal": @"Normal",
		@"speedFast": @"Fast",
		@"petCount": @"Visible pets",
		@"petCountUnit": @" pets",
		@"keyboardReaction": @"Keyboard reaction",
		@"language": @"Language",
		@"display": @"Display",
		@"quit": @"Quit",
		@"settingsTitle": @"Animals Desktop Settings",
		@"support": @"Supported OS: macOS 12 Monterey or later / Intel and Apple Silicon",
		@"tabAnimals": @"Animals",
		@"tabMotion": @"Motion",
		@"tabNames": @"Names",
		@"visibleCount": @"Visible pets",
		@"animalMode": @"Animal selection",
		@"coatFixed": @"Fixed",
		@"coatSelected": @"Choose each",
		@"coatRandom": @"Random",
		@"fixedAnimal": @"Fixed animal",
		@"perPetAnimal": @"Per-pet animals",
		@"petIndexSuffix": @"",
		@"mode": @"Mode",
		@"modeKeyboard": @"Keyboard reaction",
		@"modeRandom": @"Random stroll",
		@"speedLabel": @"Speed",
		@"typingWheel": @"Chinchilla/hamster wheel",
		@"petSize": @"Size",
		@"nameLabels": @"Show names",
		@"nameHint": @"When enabled, names appear while the pointer is over an animal.",
		@"defaultPetName": @"Animal ",
		@"macNote": @"The Mac app stays in the menu bar and normally does not show a Dock icon.",
		@"close": @"Close",
	};
	NSString *value = (english ? [en objectForKey:key] : [ja objectForKey:key]);
	return value != nil ? value : key;
}

static NSString *AnimalsDisplayLabel(NSInteger index, NSScreen *screen) {
	NSRect visible = [screen visibleFrame];
	BOOL english = goAnimalsDesktopGetLanguage() == 1;
	NSString *base = english
		? [NSString stringWithFormat:@"Display %ld", (long)index + 1]
		: [NSString stringWithFormat:@"ディスプレイ %ld", (long)index + 1];
	if (screen == [NSScreen mainScreen]) {
		base = english
			? [base stringByAppendingString:@" (Main)"]
			: [base stringByAppendingString:@" (メイン)"];
	}
	return [NSString stringWithFormat:@"%@ %.0fx%.0f", base, visible.size.width, visible.size.height];
}

@interface AnimalsView : NSView
@property(nonatomic, retain) NSImage *image;
@property(nonatomic) NSInteger hoverPet;
@end

@implementation AnimalsView
- (instancetype)initWithFrame:(NSRect)frame {
	self = [super initWithFrame:frame];
	if (self) {
		_hoverPet = -1;
	}
	return self;
}

- (BOOL)isOpaque {
	return NO;
}

- (BOOL)isFlipped {
	return YES;
}

- (void)drawRect:(NSRect)dirtyRect {
	[[NSColor clearColor] setFill];
	NSRectFill(self.bounds);
	if (self.image != nil) {
		[self.image drawInRect:self.bounds];
	}
	if (goAnimalsDesktopGetNameLabels() == 0 || self.hoverPet < 0 || self.hoverPet >= AnimalsMaxPetCount) {
		return;
	}
	char buffer[256] = {0};
	int copied = goAnimalsDesktopCopyPetName((int)self.hoverPet, buffer, (int)sizeof(buffer));
	NSString *name = nil;
	if (copied > 0) {
		name = [NSString stringWithUTF8String:buffer];
	}
	if (name == nil || [name length] == 0) {
		name = [NSString stringWithFormat:@"%@%ld", AnimalsText(@"defaultPetName"), (long)self.hoverPet + 1];
	}

	NSMutableParagraphStyle *style = [[[NSMutableParagraphStyle alloc] init] autorelease];
	[style setAlignment:NSTextAlignmentCenter];
	[style setLineBreakMode:NSLineBreakByTruncatingTail];
	NSDictionary *attrs = @{
		NSFontAttributeName: [NSFont systemFontOfSize:11.0 weight:NSFontWeightSemibold],
		NSForegroundColorAttributeName: [NSColor colorWithCalibratedWhite:1.0 alpha:0.96],
		NSParagraphStyleAttributeName: style
	};
	NSSize textSize = [name boundingRectWithSize:NSMakeSize(200.0, 20.0)
	                                    options:NSStringDrawingUsesLineFragmentOrigin
	                                 attributes:attrs].size;
	CGFloat labelW = MIN(MAX(72.0, ceil(textSize.width) + 22.0), 220.0);
	CGFloat labelH = 24.0;
	CGFloat petX = (CGFloat)goAnimalsDesktopGetPetDrawX((int)self.hoverPet);
	CGFloat petY = (CGFloat)goAnimalsDesktopGetPetDrawY((int)self.hoverPet);
	CGFloat petW = AnimalsSpriteWidth * (CGFloat)goAnimalsDesktopGetPetSize((int)self.hoverPet) / 100.0;
	CGFloat x = MIN(MAX(2.0, petX + petW / 2.0 - labelW / 2.0), MAX(2.0, self.bounds.size.width - labelW - 2.0));
	CGFloat y = MAX(0.0, petY - labelH - 4.0);
	NSRect labelRect = NSMakeRect(x, y, labelW, labelH);
	NSBezierPath *path = [NSBezierPath bezierPathWithRoundedRect:labelRect xRadius:9.0 yRadius:9.0];
	[[NSColor colorWithCalibratedRed:0.13 green:0.18 blue:0.15 alpha:0.78] setFill];
	[path fill];
	[name drawInRect:NSInsetRect(labelRect, 10.0, 4.0) withAttributes:attrs];
}
@end

@interface AnimalsAppDelegate : NSObject <NSApplicationDelegate, NSMenuDelegate, NSTextFieldDelegate>
@property(nonatomic) CGFloat sceneHeight;
@property(nonatomic, retain) NSWindow *window;
@property(nonatomic, retain) AnimalsView *view;
@property(nonatomic, retain) NSStatusItem *statusItem;
@property(nonatomic, retain) NSImage *statusIcon;
@property(nonatomic, retain) NSTimer *timer;
@property(nonatomic, retain) id globalMonitor;
@property(nonatomic, retain) id localMonitor;
@property(nonatomic, retain) id mouseClickMonitor;
@property(nonatomic, retain) id mouseMoveMonitor;
@property(nonatomic, retain) NSWindow *settingsWindow;
@property(nonatomic, retain) NSPopUpButton *countPopup;
@property(nonatomic, retain) NSPopUpButton *languagePopup;
@property(nonatomic, retain) NSPopUpButton *displayPopup;
@property(nonatomic, retain) NSPopUpButton *modePopup;
@property(nonatomic, retain) NSPopUpButton *speedPopup;
@property(nonatomic, retain) NSPopUpButton *coatModePopup;
@property(nonatomic, retain) NSPopUpButton *fixedCoatPopup;
@property(nonatomic, retain) NSMutableArray *selectedCoatPopups;
@property(nonatomic, retain) NSMutableArray *petNameFields;
@property(nonatomic, retain) NSMutableArray *petSizePopups;
@property(nonatomic, retain) NSButton *wheelCheckbox;
@property(nonatomic, retain) NSButton *nameLabelsCheckbox;
- (instancetype)initWithSceneHeight:(CGFloat)sceneHeight iconBytes:(const unsigned char *)iconBytes iconLength:(int)iconLength;
- (void)addGroupedVariantItemsToMenu:(NSMenu *)menu;
@end

static AnimalsAppDelegate *animalsDelegate = nil;

@implementation AnimalsAppDelegate
- (instancetype)initWithSceneHeight:(CGFloat)sceneHeight iconBytes:(const unsigned char *)iconBytes iconLength:(int)iconLength {
	self = [super init];
	if (self) {
		_sceneHeight = sceneHeight;
		if (iconBytes != NULL && iconLength > 0) {
			NSData *data = [NSData dataWithBytes:iconBytes length:(NSUInteger)iconLength];
			_statusIcon = [[NSImage alloc] initWithData:data];
			[_statusIcon setSize:NSMakeSize(20, 18)];
			[_statusIcon setTemplate:NO];
		}
	}
	return self;
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
	NSScreen *screen = AnimalsScreenForDisplayID(goAnimalsDesktopGetDisplayID());
	NSRect visible = [screen visibleFrame];
	CGFloat width = visible.size.width;
	if (width < 320.0) {
		width = 320.0;
	}

	NSRect frame = NSMakeRect(visible.origin.x, visible.origin.y, width, self.sceneHeight);
	self.window = [[[NSWindow alloc] initWithContentRect:frame
	                                           styleMask:NSWindowStyleMaskBorderless
	                                             backing:NSBackingStoreBuffered
	                                               defer:NO] autorelease];
	[self.window setReleasedWhenClosed:NO];
	[self.window setOpaque:NO];
	[self.window setBackgroundColor:[NSColor clearColor]];
	[self.window setHasShadow:NO];
	[self.window setIgnoresMouseEvents:YES];
	[self.window setCanHide:NO];
	[self.window setLevel:NSStatusWindowLevel];
	[self.window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces |
	                                    NSWindowCollectionBehaviorStationary |
	                                    NSWindowCollectionBehaviorIgnoresCycle];

	self.view = [[[AnimalsView alloc] initWithFrame:NSMakeRect(0, 0, width, self.sceneHeight)] autorelease];
	[self.view setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
	[self.window setContentView:self.view];
	[self.window orderFrontRegardless];
	goAnimalsDesktopSetSceneWidth((int)width);
	goAnimalsDesktopSetDisplayID(AnimalsScreenID(screen));

	self.statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
	if (self.statusIcon != nil) {
		self.statusItem.button.image = self.statusIcon;
		self.statusItem.button.imagePosition = NSImageOnly;
	} else {
		self.statusItem.button.title = @"Animals";
	}
	self.statusItem.button.toolTip = @"Animals Desktop";

	[self installStatusMenu];
	[self refreshMenuState];

	self.timer = [NSTimer timerWithTimeInterval:0.016
	                                     target:self
	                                   selector:@selector(tick:)
	                                   userInfo:nil
	                                    repeats:YES];
	[[NSRunLoop mainRunLoop] addTimer:self.timer forMode:NSRunLoopCommonModes];

	self.globalMonitor = [NSEvent addGlobalMonitorForEventsMatchingMask:NSEventMaskKeyDown
	                                                             handler:^(NSEvent *event) {
		goAnimalsDesktopKeyDown();
	}];
	self.localMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskKeyDown
	                                                          handler:^NSEvent *(NSEvent *event) {
		goAnimalsDesktopKeyDown();
		return event;
	}];
	self.mouseClickMonitor = [NSEvent addGlobalMonitorForEventsMatchingMask:NSEventMaskLeftMouseDown
	                                                                handler:^(NSEvent *event) {
		[self handleGlobalMouseDown:event];
	}];
	self.mouseMoveMonitor = [NSEvent addGlobalMonitorForEventsMatchingMask:NSEventMaskMouseMoved
	                                                               handler:^(NSEvent *event) {
		[self updateHoverFromMouseLocation:[NSEvent mouseLocation]];
	}];
	[[NSNotificationCenter defaultCenter] addObserver:self
	                                         selector:@selector(screenParametersChanged:)
	                                             name:NSApplicationDidChangeScreenParametersNotification
	                                           object:nil];
}

- (void)installStatusMenu {
	NSMenu *menu = [[[NSMenu alloc] initWithTitle:@"Animals Desktop"] autorelease];
	[menu setDelegate:self];
	NSMenuItem *title = [[[NSMenuItem alloc] initWithTitle:@"Animals Desktop" action:nil keyEquivalent:@""] autorelease];
	[title setEnabled:NO];
	[menu addItem:title];
	[menu addItem:[NSMenuItem separatorItem]];

	NSMenuItem *settings = [self menuItemWithTitle:AnimalsText(@"settingsOpen") action:@selector(showSettings:) tag:AnimalsMenuSettings];
	[menu addItem:settings];
	[menu addItem:[NSMenuItem separatorItem]];

	NSMenu *languageMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"language")] autorelease];
	[languageMenu setDelegate:self];
	[languageMenu addItem:[self menuItemWithTitle:@"日本語" action:@selector(setLanguage:) tag:AnimalsMenuLanguageJA]];
	[languageMenu addItem:[self menuItemWithTitle:@"English" action:@selector(setLanguage:) tag:AnimalsMenuLanguageEN]];
	NSMenuItem *languageRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"language") action:nil keyEquivalent:@""] autorelease];
	[languageRoot setSubmenu:languageMenu];
	[menu addItem:languageRoot];

	NSMenu *displayMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"display")] autorelease];
	[displayMenu setDelegate:self];
	NSArray *screens = [NSScreen screens];
	for (NSInteger i = 0; i < [screens count]; i++) {
		NSScreen *screen = [screens objectAtIndex:i];
		[displayMenu addItem:[self menuItemWithTitle:AnimalsDisplayLabel(i, screen) action:@selector(setDisplayFromMenu:) tag:AnimalsMenuDisplayBase + i]];
	}
	NSMenuItem *displayRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"display") action:nil keyEquivalent:@""] autorelease];
	[displayRoot setSubmenu:displayMenu];
	[menu addItem:displayRoot];

	NSMenu *variantMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"fixedAnimal")] autorelease];
	[variantMenu setDelegate:self];
	[self addGroupedVariantItemsToMenu:variantMenu];
	NSMenuItem *variantRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"fixedAnimal") action:nil keyEquivalent:@""] autorelease];
	[variantRoot setSubmenu:variantMenu];
	[menu addItem:variantRoot];

	NSMenu *coatModeMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"animalMode")] autorelease];
	[coatModeMenu setDelegate:self];
	[coatModeMenu addItem:[self menuItemWithTitle:AnimalsText(@"coatFixed") action:@selector(setCoatModeFromMenu:) tag:AnimalsMenuCoatFixed]];
	[coatModeMenu addItem:[self menuItemWithTitle:AnimalsText(@"coatSelected") action:@selector(setCoatModeFromMenu:) tag:AnimalsMenuCoatSelected]];
	[coatModeMenu addItem:[self menuItemWithTitle:AnimalsText(@"coatRandom") action:@selector(setCoatModeFromMenu:) tag:AnimalsMenuCoatRandom]];
	NSMenuItem *coatModeRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"animalMode") action:nil keyEquivalent:@""] autorelease];
	[coatModeRoot setSubmenu:coatModeMenu];
	[menu addItem:coatModeRoot];

	NSMenu *modeMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"mode")] autorelease];
	[modeMenu setDelegate:self];
	[modeMenu addItem:[self menuItemWithTitle:AnimalsText(@"modeKeyboard") action:@selector(setModeFromMenu:) tag:AnimalsMenuModeKeyboard]];
	[modeMenu addItem:[self menuItemWithTitle:AnimalsText(@"modeRandom") action:@selector(setModeFromMenu:) tag:AnimalsMenuModeRandom]];
	NSMenuItem *modeRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"mode") action:nil keyEquivalent:@""] autorelease];
	[modeRoot setSubmenu:modeMenu];
	[menu addItem:modeRoot];

	NSMenu *speedMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"speed")] autorelease];
	[speedMenu setDelegate:self];
	[speedMenu addItem:[self menuItemWithTitle:AnimalsText(@"speedSlow") action:@selector(setSpeed:) tag:AnimalsMenuSpeedSlow]];
	[speedMenu addItem:[self menuItemWithTitle:AnimalsText(@"speedNormal") action:@selector(setSpeed:) tag:AnimalsMenuSpeedNormal]];
	[speedMenu addItem:[self menuItemWithTitle:AnimalsText(@"speedFast") action:@selector(setSpeed:) tag:AnimalsMenuSpeedFast]];
	NSMenuItem *speedRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"speed") action:nil keyEquivalent:@""] autorelease];
	[speedRoot setSubmenu:speedMenu];
	[menu addItem:speedRoot];

	NSMenu *countMenu = [[[NSMenu alloc] initWithTitle:AnimalsText(@"petCount")] autorelease];
	[countMenu setDelegate:self];
	for (NSInteger i = 1; i <= AnimalsMaxPetCount; i++) {
		NSString *label = [NSString stringWithFormat:@"%ld%@", (long)i, AnimalsText(@"petCountUnit")];
		[countMenu addItem:[self menuItemWithTitle:label action:@selector(setPetCount:) tag:AnimalsMenuCountBase + i]];
	}
	NSMenuItem *countRoot = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"petCount") action:nil keyEquivalent:@""] autorelease];
	[countRoot setSubmenu:countMenu];
	[menu addItem:countRoot];

	NSMenuItem *wheel = [self menuItemWithTitle:AnimalsText(@"typingWheel") action:@selector(toggleWheelEnabled:) tag:AnimalsMenuWheelEnabled];
	[menu addItem:wheel];
	[menu addItem:[NSMenuItem separatorItem]];

	NSMenuItem *quit = [[[NSMenuItem alloc] initWithTitle:AnimalsText(@"quit") action:@selector(quit:) keyEquivalent:@"q"] autorelease];
	[quit setTarget:self];
	[menu addItem:quit];
	self.statusItem.menu = menu;
}

- (void)applicationWillTerminate:(NSNotification *)notification {
	if (self.globalMonitor != nil) {
		[NSEvent removeMonitor:self.globalMonitor];
		self.globalMonitor = nil;
	}
	if (self.localMonitor != nil) {
		[NSEvent removeMonitor:self.localMonitor];
		self.localMonitor = nil;
	}
	if (self.mouseClickMonitor != nil) {
		[NSEvent removeMonitor:self.mouseClickMonitor];
		self.mouseClickMonitor = nil;
	}
	if (self.mouseMoveMonitor != nil) {
		[NSEvent removeMonitor:self.mouseMoveMonitor];
		self.mouseMoveMonitor = nil;
	}
	[[NSNotificationCenter defaultCenter] removeObserver:self];
	[self.timer invalidate];
}

- (void)tick:(NSTimer *)timer {
	goAnimalsDesktopTick();
}

- (NSMenuItem *)menuItemWithTitle:(NSString *)title action:(SEL)action tag:(NSInteger)tag {
	NSMenuItem *item = [[[NSMenuItem alloc] initWithTitle:title action:action keyEquivalent:@""] autorelease];
	[item setTarget:self];
	[item setTag:tag];
	return item;
}

- (void)addGroupedVariantItemsToMenu:(NSMenu *)menu {
	NSMutableDictionary *menusByGroup = [NSMutableDictionary dictionary];
	for (NSInteger i = 0; i < goAnimalsDesktopGetVariantCount(); i++) {
		NSString *group = AnimalsVariantGroupLabel(i);
		NSMenu *groupMenu = [menusByGroup objectForKey:group];
		if (groupMenu == nil) {
			groupMenu = [[[NSMenu alloc] initWithTitle:group] autorelease];
			[groupMenu setDelegate:self];
			NSMenuItem *root = [[[NSMenuItem alloc] initWithTitle:group action:nil keyEquivalent:@""] autorelease];
			[root setSubmenu:groupMenu];
			[menu addItem:root];
			[menusByGroup setObject:groupMenu forKey:group];
		}
		[groupMenu addItem:[self menuItemWithTitle:AnimalsVariantLabel(i) action:@selector(setVariantFromMenu:) tag:AnimalsMenuVariantBase + i]];
	}
}

- (void)menuNeedsUpdate:(NSMenu *)menu {
	[self refreshMenuState];
}

- (void)menuWillOpen:(NSMenu *)menu {
	[self refreshMenuState];
}

- (void)refreshMenuState {
	if (self.statusItem == nil || self.statusItem.menu == nil) {
		return;
	}
	[self refreshMenuState:self.statusItem.menu speed:goAnimalsDesktopGetSpeed() count:goAnimalsDesktopGetPetCount() wheelEnabled:goAnimalsDesktopGetWheelEnabled() language:goAnimalsDesktopGetLanguage() displayID:AnimalsActiveDisplayID() coatMode:goAnimalsDesktopGetCoatMode() variant:goAnimalsDesktopGetVariant() mode:goAnimalsDesktopGetMode()];
	[self refreshSettingsControls];
}

- (void)refreshMenuState:(NSMenu *)menu speed:(int)speed count:(int)count wheelEnabled:(int)wheelEnabled language:(int)language displayID:(long long)displayID coatMode:(int)coatMode variant:(int)variant mode:(int)mode {
	NSArray *screens = [NSScreen screens];
	for (NSMenuItem *item in [menu itemArray]) {
		NSInteger tag = [item tag];
		if (tag == AnimalsMenuSpeedSlow || tag == AnimalsMenuSpeedNormal || tag == AnimalsMenuSpeedFast) {
			int itemSpeed = 3;
			if (tag == AnimalsMenuSpeedSlow) {
				itemSpeed = 2;
			} else if (tag == AnimalsMenuSpeedFast) {
				itemSpeed = 5;
			}
			[item setState:(itemSpeed == speed) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag > AnimalsMenuCountBase && tag <= AnimalsMenuCountBase + AnimalsMaxPetCount) {
			int itemCount = (int)(tag - AnimalsMenuCountBase);
			[item setState:(itemCount == count) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuWheelEnabled) {
			[item setState:wheelEnabled ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuLanguageJA) {
			[item setState:(language == 0) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuLanguageEN) {
			[item setState:(language == 1) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag >= AnimalsMenuDisplayBase && tag < AnimalsMenuDisplayBase + [screens count]) {
			NSScreen *screen = [screens objectAtIndex:(tag - AnimalsMenuDisplayBase)];
			[item setState:(AnimalsScreenID(screen) == displayID) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuCoatFixed) {
			[item setState:(coatMode == 0) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuCoatSelected) {
			[item setState:(coatMode == 1) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuCoatRandom) {
			[item setState:(coatMode == 2) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuModeKeyboard) {
			[item setState:(mode == 0) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag == AnimalsMenuModeRandom) {
			[item setState:(mode == 1) ? NSControlStateValueOn : NSControlStateValueOff];
		} else if (tag >= AnimalsMenuVariantBase && tag < AnimalsMenuVariantBase + goAnimalsDesktopGetVariantCount()) {
			[item setState:((int)(tag - AnimalsMenuVariantBase) == variant) ? NSControlStateValueOn : NSControlStateValueOff];
		}
		if ([item submenu] != nil) {
			[self refreshMenuState:[item submenu] speed:speed count:count wheelEnabled:wheelEnabled language:language displayID:displayID coatMode:coatMode variant:variant mode:mode];
		}
	}
}

- (void)applyDisplaySelection {
	if (self.window == nil) {
		return;
	}
	NSScreen *screen = AnimalsScreenForDisplayID(goAnimalsDesktopGetDisplayID());
	long long displayID = AnimalsScreenID(screen);
	if (displayID != goAnimalsDesktopGetDisplayID()) {
		goAnimalsDesktopSetDisplayID(displayID);
	}
	NSRect visible = [screen visibleFrame];
	CGFloat width = visible.size.width;
	if (width < 320.0) {
		width = 320.0;
	}
	NSRect frame = NSMakeRect(visible.origin.x, visible.origin.y, width, self.sceneHeight);
	[self.window setFrame:frame display:YES];
	[self.view setFrame:NSMakeRect(0, 0, width, self.sceneHeight)];
	goAnimalsDesktopSetSceneWidth((int)width);
	[self updateHoverFromMouseLocation:[NSEvent mouseLocation]];
}

- (void)screenParametersChanged:(NSNotification *)notification {
	(void)notification;
	[self applyDisplaySelection];
	[self installStatusMenu];
	if (self.settingsWindow != nil) {
		BOOL visible = [self.settingsWindow isVisible];
		[self.settingsWindow close];
		self.settingsWindow = nil;
		if (visible) {
			[self ensureSettingsWindow];
			[self.settingsWindow makeKeyAndOrderFront:nil];
		}
	}
	[self refreshMenuState];
}

- (void)setLanguage:(id)sender {
	NSInteger tag = [sender tag];
	goAnimalsDesktopSetLanguage(tag == AnimalsMenuLanguageEN ? 1 : 0);
	[self installStatusMenu];
	if (self.settingsWindow != nil) {
		[self.settingsWindow close];
		self.settingsWindow = nil;
		[self ensureSettingsWindow];
		[self.settingsWindow makeKeyAndOrderFront:nil];
	}
	[self refreshMenuState];
}

- (void)setDisplayFromMenu:(id)sender {
	NSInteger screenIndex = [sender tag] - AnimalsMenuDisplayBase;
	NSArray *screens = [NSScreen screens];
	if (screenIndex < 0 || screenIndex >= [screens count]) {
		return;
	}
	NSScreen *screen = [screens objectAtIndex:screenIndex];
	goAnimalsDesktopSetDisplayID(AnimalsScreenID(screen));
	[self applyDisplaySelection];
	[self refreshMenuState];
}

- (void)setVariantFromMenu:(id)sender {
	goAnimalsDesktopSetVariant((int)([sender tag] - AnimalsMenuVariantBase));
	goAnimalsDesktopSetCoatMode(0);
	[self refreshMenuState];
}

- (void)setCoatModeFromMenu:(id)sender {
	NSInteger tag = [sender tag];
	if (tag == AnimalsMenuCoatFixed) {
		goAnimalsDesktopSetCoatMode(0);
	} else if (tag == AnimalsMenuCoatSelected) {
		goAnimalsDesktopSetCoatMode(1);
	} else {
		goAnimalsDesktopSetCoatMode(2);
	}
	[self refreshMenuState];
}

- (void)setModeFromMenu:(id)sender {
	goAnimalsDesktopSetMode(([sender tag] == AnimalsMenuModeKeyboard) ? 0 : 1);
	[self refreshMenuState];
}

- (void)setSpeed:(id)sender {
	NSInteger tag = [sender tag];
	if (tag == AnimalsMenuSpeedSlow) {
		goAnimalsDesktopSetSpeed(2);
	} else if (tag == AnimalsMenuSpeedFast) {
		goAnimalsDesktopSetSpeed(5);
	} else {
		goAnimalsDesktopSetSpeed(3);
	}
	[self refreshMenuState];
}

- (void)setPetCount:(id)sender {
	NSInteger tag = [sender tag];
	if (tag > AnimalsMenuCountBase && tag <= AnimalsMenuCountBase + AnimalsMaxPetCount) {
		goAnimalsDesktopSetPetCount((int)(tag - AnimalsMenuCountBase));
	}
	[self refreshMenuState];
}

- (void)toggleWheelEnabled:(id)sender {
	goAnimalsDesktopSetWheelEnabled(goAnimalsDesktopGetWheelEnabled() ? 0 : 1);
	[self refreshMenuState];
}

- (void)showSettings:(id)sender {
	[self ensureSettingsWindow];
	[self refreshSettingsControls];
	[self.settingsWindow makeKeyAndOrderFront:nil];
	[NSApp activateIgnoringOtherApps:YES];
}

- (NSTextField *)labelWithTitle:(NSString *)title frame:(NSRect)frame {
	NSTextField *label = [[[NSTextField alloc] initWithFrame:frame] autorelease];
	[label setStringValue:title];
	[label setBezeled:NO];
	[label setDrawsBackground:NO];
	[label setEditable:NO];
	[label setSelectable:NO];
	[label setFont:[NSFont systemFontOfSize:12.0]];
	return label;
}

- (NSPopUpButton *)popupWithFrame:(NSRect)frame action:(SEL)action {
	NSPopUpButton *popup = [[[NSPopUpButton alloc] initWithFrame:frame pullsDown:NO] autorelease];
	[popup setTarget:self];
	[popup setAction:action];
	return popup;
}

- (void)handleGlobalMouseDown:(NSEvent *)event {
	(void)event;
	NSPoint point = [NSEvent mouseLocation];
	int localX = 0;
	int localY = 0;
	if (![self localScenePointFromScreenPoint:point x:&localX y:&localY]) {
		return;
	}
	goAnimalsDesktopClick(localX, localY);
	[self updateHoverFromSceneX:localX y:localY];
}

- (BOOL)localScenePointFromScreenPoint:(NSPoint)point x:(int *)x y:(int *)y {
	if (self.window == nil) {
		return NO;
	}
	NSRect frame = [self.window frame];
	if (!NSPointInRect(point, frame)) {
		return NO;
	}
	CGFloat localX = point.x - frame.origin.x;
	CGFloat localY = self.sceneHeight - (point.y - frame.origin.y);
	if (localX < 0.0 || localY < 0.0 || localX >= frame.size.width || localY >= self.sceneHeight) {
		return NO;
	}
	if (x != NULL) {
		*x = (int)floor(localX);
	}
	if (y != NULL) {
		*y = (int)floor(localY);
	}
	return YES;
}

- (void)updateHoverFromMouseLocation:(NSPoint)point {
	int localX = 0;
	int localY = 0;
	if (![self localScenePointFromScreenPoint:point x:&localX y:&localY]) {
		[self updateHoverPet:-1];
		return;
	}
	[self updateHoverFromSceneX:localX y:localY];
}

- (void)updateHoverFromSceneX:(int)localX y:(int)localY {
	if (goAnimalsDesktopGetNameLabels() == 0) {
		[self updateHoverPet:-1];
		return;
	}
	[self updateHoverPet:goAnimalsDesktopPetAt(localX, localY)];
}

- (void)updateHoverPet:(NSInteger)index {
	if (self.view == nil || self.view.hoverPet == index) {
		return;
	}
	self.view.hoverPet = index;
	[self.view setNeedsDisplay:YES];
}

- (void)ensureSettingsWindow {
	if (self.settingsWindow != nil) {
		return;
	}
	NSRect frame = NSMakeRect(0, 0, 620, 560);
	self.settingsWindow = [[[NSWindow alloc] initWithContentRect:frame
	                                                   styleMask:NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable
	                                                     backing:NSBackingStoreBuffered
	                                                       defer:NO] autorelease];
	[self.settingsWindow setTitle:AnimalsText(@"settingsTitle")];
	[self.settingsWindow setReleasedWhenClosed:NO];
	[self.settingsWindow center];

	NSView *content = [[[NSView alloc] initWithFrame:frame] autorelease];
	[content setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
	[self.settingsWindow setContentView:content];

	NSTextField *title = [self labelWithTitle:AnimalsText(@"settingsTitle") frame:NSMakeRect(24, 516, 360, 24)];
	[title setFont:[NSFont boldSystemFontOfSize:18.0]];
	[content addSubview:title];

	NSTextField *support = [self labelWithTitle:AnimalsText(@"support") frame:NSMakeRect(24, 492, 540, 20)];
	[support setTextColor:[NSColor secondaryLabelColor]];
	[content addSubview:support];

	NSTabView *tabs = [[[NSTabView alloc] initWithFrame:NSMakeRect(20, 58, 580, 420)] autorelease];
	[content addSubview:tabs];

	NSTabViewItem *animals = [[[NSTabViewItem alloc] initWithIdentifier:@"animals"] autorelease];
	[animals setLabel:AnimalsText(@"tabAnimals")];
	NSView *animalView = [[[NSView alloc] initWithFrame:NSMakeRect(0, 0, 580, 392)] autorelease];
	[animals setView:animalView];
	[tabs addTabViewItem:animals];

	NSTabViewItem *motion = [[[NSTabViewItem alloc] initWithIdentifier:@"motion"] autorelease];
	[motion setLabel:AnimalsText(@"tabMotion")];
	NSView *motionView = [[[NSView alloc] initWithFrame:NSMakeRect(0, 0, 580, 392)] autorelease];
	[motion setView:motionView];
	[tabs addTabViewItem:motion];

	NSTabViewItem *names = [[[NSTabViewItem alloc] initWithIdentifier:@"names"] autorelease];
	[names setLabel:AnimalsText(@"tabNames")];
	NSView *namesView = [[[NSView alloc] initWithFrame:NSMakeRect(0, 0, 580, 392)] autorelease];
	[names setView:namesView];
	[tabs addTabViewItem:names];

	[animalView addSubview:[self labelWithTitle:AnimalsText(@"visibleCount") frame:NSMakeRect(22, 340, 120, 24)]];
	self.countPopup = [self popupWithFrame:NSMakeRect(150, 336, 180, 28) action:@selector(settingsCountChanged:)];
	for (NSInteger i = 1; i <= AnimalsMaxPetCount; i++) {
		[self.countPopup addItemWithTitle:[NSString stringWithFormat:@"%ld%@", (long)i, AnimalsText(@"petCountUnit")]];
	}
	[animalView addSubview:self.countPopup];

	[animalView addSubview:[self labelWithTitle:AnimalsText(@"language") frame:NSMakeRect(22, 298, 120, 24)]];
	self.languagePopup = [self popupWithFrame:NSMakeRect(150, 294, 180, 28) action:@selector(settingsLanguageChanged:)];
	for (NSString *label in @[@"日本語", @"English"]) {
		[self.languagePopup addItemWithTitle:label];
	}
	[animalView addSubview:self.languagePopup];

	[animalView addSubview:[self labelWithTitle:AnimalsText(@"animalMode") frame:NSMakeRect(22, 256, 120, 24)]];
	self.coatModePopup = [self popupWithFrame:NSMakeRect(150, 294, 180, 28) action:@selector(settingsCoatModeChanged:)];
	[self.coatModePopup setFrame:NSMakeRect(150, 252, 180, 28)];
	for (NSString *label in @[AnimalsText(@"coatFixed"), AnimalsText(@"coatSelected"), AnimalsText(@"coatRandom")]) {
		[self.coatModePopup addItemWithTitle:label];
	}
	[animalView addSubview:self.coatModePopup];

	[animalView addSubview:[self labelWithTitle:AnimalsText(@"fixedAnimal") frame:NSMakeRect(22, 214, 120, 24)]];
	self.fixedCoatPopup = [self popupWithFrame:NSMakeRect(150, 252, 260, 28) action:@selector(settingsFixedCoatChanged:)];
	[self.fixedCoatPopup setFrame:NSMakeRect(150, 210, 260, 28)];
	[self populateVariantPopup:self.fixedCoatPopup];
	[animalView addSubview:self.fixedCoatPopup];

	NSTextField *perPet = [self labelWithTitle:AnimalsText(@"perPetAnimal") frame:NSMakeRect(22, 172, 160, 24)];
	[perPet setFont:[NSFont boldSystemFontOfSize:12.0]];
	[animalView addSubview:perPet];

	self.selectedCoatPopups = [NSMutableArray arrayWithCapacity:AnimalsMaxPetCount];
	for (NSInteger i = 0; i < AnimalsMaxPetCount; i++) {
		NSInteger column = i / 5;
		NSInteger row = i % 5;
		CGFloat x = 22 + column * 270;
		CGFloat y = 132 - row * 30;
		NSString *petLabel = (goAnimalsDesktopGetLanguage() == 1)
			? [NSString stringWithFormat:@"%ld", (long)i + 1]
			: [NSString stringWithFormat:@"%ld%@", (long)i + 1, AnimalsText(@"petIndexSuffix")];
		[animalView addSubview:[self labelWithTitle:petLabel frame:NSMakeRect(x, y + 4, 58, 22)]];
		NSPopUpButton *popup = [self popupWithFrame:NSMakeRect(x + 62, y, 190, 28) action:@selector(settingsSelectedCoatChanged:)];
		[popup setTag:i];
		[self populateVariantPopup:popup];
		[animalView addSubview:popup];
		[self.selectedCoatPopups addObject:popup];
	}

	[motionView addSubview:[self labelWithTitle:AnimalsText(@"mode") frame:NSMakeRect(22, 340, 120, 24)]];
	self.modePopup = [self popupWithFrame:NSMakeRect(150, 336, 220, 28) action:@selector(settingsModeChanged:)];
	for (NSString *label in @[AnimalsText(@"modeKeyboard"), AnimalsText(@"modeRandom")]) {
		[self.modePopup addItemWithTitle:label];
	}
	[motionView addSubview:self.modePopup];

	[motionView addSubview:[self labelWithTitle:AnimalsText(@"speedLabel") frame:NSMakeRect(22, 298, 120, 24)]];
	self.speedPopup = [self popupWithFrame:NSMakeRect(150, 294, 180, 28) action:@selector(settingsSpeedChanged:)];
	for (NSString *label in @[AnimalsText(@"speedSlow"), AnimalsText(@"speedNormal"), AnimalsText(@"speedFast")]) {
		[self.speedPopup addItemWithTitle:label];
	}
	[motionView addSubview:self.speedPopup];

	[motionView addSubview:[self labelWithTitle:AnimalsText(@"display") frame:NSMakeRect(22, 256, 120, 24)]];
	self.displayPopup = [self popupWithFrame:NSMakeRect(150, 252, 260, 28) action:@selector(settingsDisplayChanged:)];
	[self populateDisplayPopup:self.displayPopup];
	[motionView addSubview:self.displayPopup];

	self.wheelCheckbox = [[[NSButton alloc] initWithFrame:NSMakeRect(150, 248, 260, 28)] autorelease];
	[self.wheelCheckbox setFrame:NSMakeRect(150, 206, 260, 28)];
	[self.wheelCheckbox setButtonType:NSButtonTypeSwitch];
	[self.wheelCheckbox setTitle:AnimalsText(@"typingWheel")];
	[self.wheelCheckbox setTarget:self];
	[self.wheelCheckbox setAction:@selector(settingsWheelChanged:)];
	[motionView addSubview:self.wheelCheckbox];

	self.nameLabelsCheckbox = [[[NSButton alloc] initWithFrame:NSMakeRect(22, 340, 240, 28)] autorelease];
	[self.nameLabelsCheckbox setButtonType:NSButtonTypeSwitch];
	[self.nameLabelsCheckbox setTitle:AnimalsText(@"nameLabels")];
	[self.nameLabelsCheckbox setTarget:self];
	[self.nameLabelsCheckbox setAction:@selector(settingsNameLabelsChanged:)];
	[namesView addSubview:self.nameLabelsCheckbox];

	NSTextField *nameHint = [self labelWithTitle:AnimalsText(@"nameHint") frame:NSMakeRect(22, 312, 520, 22)];
	[nameHint setTextColor:[NSColor secondaryLabelColor]];
	[namesView addSubview:nameHint];

	self.petNameFields = [NSMutableArray arrayWithCapacity:AnimalsMaxPetCount];
	self.petSizePopups = [NSMutableArray arrayWithCapacity:AnimalsMaxPetCount];
	[namesView addSubview:[self labelWithTitle:AnimalsText(@"petSize") frame:NSMakeRect(216, 286, 72, 22)]];
	[namesView addSubview:[self labelWithTitle:AnimalsText(@"petSize") frame:NSMakeRect(486, 286, 72, 22)]];
	for (NSInteger i = 0; i < AnimalsMaxPetCount; i++) {
		NSInteger column = i / 5;
		NSInteger row = i % 5;
		CGFloat x = 22 + column * 270;
		CGFloat y = 262 - row * 42;
		NSString *petLabel = (goAnimalsDesktopGetLanguage() == 1)
			? [NSString stringWithFormat:@"%ld", (long)i + 1]
			: [NSString stringWithFormat:@"%ld%@", (long)i + 1, AnimalsText(@"petIndexSuffix")];
		[namesView addSubview:[self labelWithTitle:petLabel frame:NSMakeRect(x, y + 4, 58, 24)]];
		NSTextField *field = [[[NSTextField alloc] initWithFrame:NSMakeRect(x + 62, y, 190, 28)] autorelease];
		[field setTag:i];
		[field setTarget:self];
		[field setAction:@selector(settingsPetNameChanged:)];
		[field setDelegate:self];
		[field setFrame:NSMakeRect(x + 62, y, 126, 28)];
		[field setPlaceholderString:[NSString stringWithFormat:@"%@%ld", AnimalsText(@"defaultPetName"), (long)i + 1]];
		[namesView addSubview:field];
		[self.petNameFields addObject:field];

		NSPopUpButton *sizePopup = [self popupWithFrame:NSMakeRect(x + 194, y, 72, 28) action:@selector(settingsPetSizeChanged:)];
		[sizePopup setTag:i];
		[self populatePetSizePopup:sizePopup];
		[namesView addSubview:sizePopup];
		[self.petSizePopups addObject:sizePopup];
	}

	NSTextField *note = [self labelWithTitle:AnimalsText(@"macNote") frame:NSMakeRect(24, 28, 480, 20)];
	[note setTextColor:[NSColor secondaryLabelColor]];
	[content addSubview:note];

	NSButton *close = [[[NSButton alloc] initWithFrame:NSMakeRect(492, 20, 92, 32)] autorelease];
	[close setTitle:AnimalsText(@"close")];
	[close setBezelStyle:NSBezelStyleRounded];
	[close setTarget:self];
	[close setAction:@selector(closeSettings:)];
	[content addSubview:close];
}

- (void)populateVariantPopup:(NSPopUpButton *)popup {
	[popup removeAllItems];
	NSMutableArray *groupOrder = [NSMutableArray array];
	NSMutableDictionary *indicesByGroup = [NSMutableDictionary dictionary];
	NSInteger count = (NSInteger)goAnimalsDesktopGetVariantCount();
	for (NSInteger i = 0; i < count; i++) {
		NSString *group = AnimalsVariantGroupLabel(i);
		NSMutableArray *indices = [indicesByGroup objectForKey:group];
		if (indices == nil) {
			indices = [NSMutableArray array];
			[indicesByGroup setObject:indices forKey:group];
			[groupOrder addObject:group];
		}
		[indices addObject:[NSNumber numberWithInteger:i]];
	}
	for (NSString *group in groupOrder) {
		[popup addItemWithTitle:group];
		[[popup lastItem] setEnabled:NO];
		[[popup lastItem] setTag:-1];
		for (NSNumber *number in [indicesByGroup objectForKey:group]) {
			NSInteger i = [number integerValue];
			[popup addItemWithTitle:AnimalsVariantDisplayLabel(i)];
			[[popup lastItem] setTag:i];
		}
	}
}

- (void)selectVariant:(NSInteger)variant inPopup:(NSPopUpButton *)popup {
	for (NSMenuItem *item in [[popup menu] itemArray]) {
		if ([item tag] == variant) {
			[popup selectItem:item];
			return;
		}
	}
}

- (void)populatePetSizePopup:(NSPopUpButton *)popup {
	[popup removeAllItems];
	for (NSInteger percent = AnimalsMinPetSize; percent <= AnimalsMaxPetSize; percent += AnimalsPetSizeStep) {
		[popup addItemWithTitle:AnimalsPetSizeLabel(percent)];
	}
}

- (void)populateDisplayPopup:(NSPopUpButton *)popup {
	[popup removeAllItems];
	NSArray *screens = [NSScreen screens];
	for (NSInteger i = 0; i < [screens count]; i++) {
		NSScreen *screen = [screens objectAtIndex:i];
		[popup addItemWithTitle:AnimalsDisplayLabel(i, screen)];
	}
}

- (void)refreshSettingsControls {
	if (self.settingsWindow == nil) {
		return;
	}
	NSInteger count = goAnimalsDesktopGetPetCount();
	[self.countPopup selectItemAtIndex:MAX(0, MIN(AnimalsMaxPetCount - 1, count - 1))];
	[self.languagePopup selectItemAtIndex:goAnimalsDesktopGetLanguage() == 1 ? 1 : 0];
	[self.coatModePopup selectItemAtIndex:goAnimalsDesktopGetCoatMode()];
	[self selectVariant:goAnimalsDesktopGetVariant() inPopup:self.fixedCoatPopup];
	[self.modePopup selectItemAtIndex:goAnimalsDesktopGetMode()];
	NSInteger speed = goAnimalsDesktopGetSpeed();
	[self.speedPopup selectItemAtIndex:(speed == 2 ? 0 : (speed == 5 ? 2 : 1))];
	NSArray *screens = [NSScreen screens];
	long long displayID = AnimalsActiveDisplayID();
	NSInteger selectedDisplayIndex = 0;
	for (NSInteger i = 0; i < [screens count]; i++) {
		if (AnimalsScreenID([screens objectAtIndex:i]) == displayID) {
			selectedDisplayIndex = i;
			break;
		}
	}
	if ([self.displayPopup numberOfItems] != [screens count]) {
		[self populateDisplayPopup:self.displayPopup];
	}
	if ([self.displayPopup numberOfItems] > 0) {
		[self.displayPopup selectItemAtIndex:MAX(0, MIN([self.displayPopup numberOfItems] - 1, selectedDisplayIndex))];
	}
	[self.wheelCheckbox setState:goAnimalsDesktopGetWheelEnabled() ? NSControlStateValueOn : NSControlStateValueOff];
	for (NSInteger i = 0; i < [self.selectedCoatPopups count]; i++) {
		NSPopUpButton *popup = [self.selectedCoatPopups objectAtIndex:i];
		[self selectVariant:goAnimalsDesktopGetSelectedCoat((int)i) inPopup:popup];
		[popup setEnabled:(goAnimalsDesktopGetCoatMode() == 1 && i < count)];
	}
	[self.fixedCoatPopup setEnabled:(goAnimalsDesktopGetCoatMode() == 0)];
	[self.nameLabelsCheckbox setState:goAnimalsDesktopGetNameLabels() ? NSControlStateValueOn : NSControlStateValueOff];
	for (NSInteger i = 0; i < [self.petNameFields count]; i++) {
		NSTextField *field = [self.petNameFields objectAtIndex:i];
		char buffer[256] = {0};
		if (goAnimalsDesktopCopyPetName((int)i, buffer, (int)sizeof(buffer)) > 0) {
			NSString *name = [NSString stringWithUTF8String:buffer];
			[field setStringValue:(name != nil ? name : @"")];
		} else {
			[field setStringValue:@""];
		}
		[field setEnabled:(goAnimalsDesktopGetNameLabels() != 0 && i < count)];
	}
	for (NSInteger i = 0; i < [self.petSizePopups count]; i++) {
		NSPopUpButton *popup = [self.petSizePopups objectAtIndex:i];
		NSInteger percent = goAnimalsDesktopGetPetSize((int)i);
		NSInteger index = (percent - AnimalsMinPetSize) / AnimalsPetSizeStep;
		[popup selectItemAtIndex:MAX(0, MIN([popup numberOfItems] - 1, index))];
		[popup setEnabled:(i < count)];
	}
}

- (void)settingsLanguageChanged:(id)sender {
	goAnimalsDesktopSetLanguage((int)[sender indexOfSelectedItem]);
	[self installStatusMenu];
	if (self.settingsWindow != nil) {
		[self.settingsWindow close];
		self.settingsWindow = nil;
		[self ensureSettingsWindow];
		[self.settingsWindow makeKeyAndOrderFront:nil];
	}
	[self refreshMenuState];
}

- (void)settingsCountChanged:(id)sender {
	NSInteger index = [sender indexOfSelectedItem];
	goAnimalsDesktopSetPetCount((int)index + 1);
	[self refreshMenuState];
}

- (void)settingsCoatModeChanged:(id)sender {
	goAnimalsDesktopSetCoatMode((int)[sender indexOfSelectedItem]);
	[self refreshMenuState];
}

- (void)settingsFixedCoatChanged:(id)sender {
	NSInteger variant = [[sender selectedItem] tag];
	if (variant >= 0) {
		goAnimalsDesktopSetVariant((int)variant);
	}
	[self refreshMenuState];
}

- (void)settingsSelectedCoatChanged:(id)sender {
	NSInteger variant = [[sender selectedItem] tag];
	if (variant >= 0) {
		goAnimalsDesktopSetSelectedCoat((int)[sender tag], (int)variant);
	}
	[self refreshMenuState];
}

- (void)settingsModeChanged:(id)sender {
	goAnimalsDesktopSetMode((int)[sender indexOfSelectedItem]);
	[self refreshMenuState];
}

- (void)settingsSpeedChanged:(id)sender {
	NSInteger index = [sender indexOfSelectedItem];
	int values[] = {2, 3, 5};
	goAnimalsDesktopSetSpeed(values[index]);
	[self refreshMenuState];
}

- (void)settingsDisplayChanged:(id)sender {
	NSArray *screens = [NSScreen screens];
	NSInteger index = [sender indexOfSelectedItem];
	if (index < 0 || index >= [screens count]) {
		return;
	}
	goAnimalsDesktopSetDisplayID(AnimalsScreenID([screens objectAtIndex:index]));
	[self applyDisplaySelection];
	[self refreshMenuState];
}

- (void)settingsWheelChanged:(id)sender {
	goAnimalsDesktopSetWheelEnabled([sender state] == NSControlStateValueOn ? 1 : 0);
	[self refreshMenuState];
}

- (void)settingsNameLabelsChanged:(id)sender {
	goAnimalsDesktopSetNameLabels([sender state] == NSControlStateValueOn ? 1 : 0);
	if (goAnimalsDesktopGetNameLabels() == 0) {
		[self updateHoverPet:-1];
	} else {
		[self updateHoverFromMouseLocation:[NSEvent mouseLocation]];
	}
	[self refreshMenuState];
}

- (void)settingsPetNameChanged:(id)sender {
	goAnimalsDesktopSetPetName((int)[sender tag], (char *)[[sender stringValue] UTF8String]);
	[self.view setNeedsDisplay:YES];
	[self refreshMenuState];
}

- (void)settingsPetSizeChanged:(id)sender {
	NSInteger percent = AnimalsMinPetSize + [sender indexOfSelectedItem] * AnimalsPetSizeStep;
	goAnimalsDesktopSetPetSize((int)[sender tag], (int)percent);
	[self.view setNeedsDisplay:YES];
	[self refreshMenuState];
}

- (void)controlTextDidEndEditing:(NSNotification *)notification {
	id object = [notification object];
	if ([object isKindOfClass:[NSTextField class]] && [self.petNameFields containsObject:object]) {
		[self settingsPetNameChanged:object];
	}
}

- (void)closeSettings:(id)sender {
	[self.settingsWindow orderOut:nil];
}

- (void)quit:(id)sender {
	[NSApp terminate:nil];
}
@end

void updateAnimalsDesktopImage(const unsigned char *bytes, int length, int width, int height) {
	if (animalsDelegate == nil || animalsDelegate.view == nil || bytes == NULL || length <= 0) {
		return;
	}
	NSData *data = [NSData dataWithBytes:bytes length:(NSUInteger)length];
	NSImage *image = [[[NSImage alloc] initWithData:data] autorelease];
	if (image == nil) {
		return;
	}
	[image setSize:NSMakeSize(width, height)];
	animalsDelegate.view.image = image;
	[animalsDelegate.view setNeedsDisplay:YES];
}

void startAnimalsDesktopApp(int sceneHeight, const unsigned char *iconBytes, int iconLength) {
	@autoreleasepool {
		NSApplication *app = [NSApplication sharedApplication];
		animalsDelegate = [[AnimalsAppDelegate alloc] initWithSceneHeight:(CGFloat)sceneHeight iconBytes:iconBytes iconLength:iconLength];
		[app setDelegate:animalsDelegate];
		[app setActivationPolicy:NSApplicationActivationPolicyAccessory];
		[app run];
	}
}
