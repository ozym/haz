select impact.add_intensity_reported('android.device.id.1' , 172.79019, -42.2, '2014-01-08T12:00:30Z'::timestamptz, 4, '');
select impact.add_intensity_reported('android.device.id.4' , 172.79019, -42.202, '2014-01-08T12:00:31Z'::timestamptz, 5, '');
select impact.add_intensity_reported('android.device.id.2' , 174.79019, -42.2, '2014-01-08T12:00:30Z'::timestamptz, 4, '');
select impact.add_intensity_reported('android.device.id.3' , 174.39019, -42.2, now(), 4, 'Wowsers');
-- insert and update
select impact.add_intensity_reported('android.device.id.5' , 172.79019, -42.2, '2014-01-08T12:00:30Z'::timestamptz, 4, 'Not so bad');
select impact.add_intensity_reported('android.device.id.5' , 172.79019, -42.2, '2014-01-08T12:00:30Z'::timestamptz, 8, 'Yikes');

select impact.add_intensity_measured('NZ.SNZO', 174.79, -40.2, '2014-01-08T12:00:30Z'::timestamptz, 1);
select impact.add_intensity_measured('NZ.WEL', 175.79, -40.2, '2014-01-08T12:00:30Z'::timestamptz, 2);
select impact.add_intensity_measured('NZ.WEL', 175.79, -40.2, '2014-01-08T12:01:30Z'::timestamptz, 5);
select impact.add_intensity_measured('NZ.FVWZ', 175.79, -39.2, now(), 5);
