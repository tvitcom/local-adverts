var agreePrivacyPolicyOnce;
var clientParams;

document.onreadystatechange = function () {
	  if (document.readyState === 'complete') {

	  	agreePrivacyPolicyOnce = function() {
	  		alert("Yo")
			when = 31536000; // Expected for year
		    document.cookie = "privacy_signer_time="+Date.now()+"; Path=/; secure; SameSite=strict; max-age="+when;
		    document.cookie = "privacy_signer_ua="+btoa(navigator.userAgent)+"; Path=/; secure; SameSite=strict; max-age="+when;
		    document.cookie = "privacy_signer_screen="+screen.availWidth+"x"+screen.availHeight+"; Path=/; SameSite=strict; secure; max-age="+when;
		    document.cookie = "privacy_signer_langs="+navigator.languages.toString()+"; Path=/; secure; SameSite=strict; max-age="+when;
			document.getElementById('modal_1').checked = false; // close modal
			console.log("Close mmodal - cookies are settings successfully")
		}
		if (!document.cookie.split(';').filter((item) => item.trim().startsWith('privacy_signer_time=')).length) {
			document.getElementById('modal_1').checked = true; // open modal
			console.log("Open mmodal - no cookies")
		}
	}
}